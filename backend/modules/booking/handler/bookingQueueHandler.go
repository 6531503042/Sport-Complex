package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/usecase"
	"main/pkg/queue"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

type (
	BookingQueueHandlerService interface {
		ProcessBookingQueue()
		HandleFacilityUpdates()
		PublishBookingCreated(booking *booking.Booking) error
	}

	bookingQueueHandler struct {
		cfg            *config.Config
		bookingUsecase usecase.BookingUsecaseService
	}
)

func NewBookingQueueHandler(cfg *config.Config, bookingUsecase usecase.BookingUsecaseService) BookingQueueHandlerService {
	return &bookingQueueHandler{
		cfg:            cfg,
		bookingUsecase: bookingUsecase,
	}
}

func (h *bookingQueueHandler) BookingConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	
	// Add more detailed logging
	log.Printf("Attempting to connect to Kafka at %s", h.cfg.Kafka.Url)
	
	// Try multiple times to connect
	var worker sarama.Consumer
	var err error
	for i := 0; i < 3; i++ {
		worker, err = sarama.NewConsumer([]string{h.cfg.Kafka.Url}, config)
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Failed to create consumer: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer after 3 attempts: %v", err)
	}

	// List topics for debugging
	topics, err := worker.Topics()
	if err != nil {
		log.Printf("Failed to get topics: %v", err)
	} else {
		log.Printf("Available topics: %v", topics)
	}

	// Create the partition consumer with retry
	var consumer sarama.PartitionConsumer
	for i := 0; i < 3; i++ {
		consumer, err = worker.ConsumePartition("bookings", 0, sarama.OffsetNewest)
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Failed to create partition consumer: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		worker.Close()
		return nil, fmt.Errorf("failed to create partition consumer after 3 attempts: %v", err)
	}

	log.Printf("Successfully connected to Kafka and created consumer")
	return consumer, nil
}

func (h *bookingQueueHandler) ProcessBookingQueue() {
	ctx := context.Background()

	consumer, err := h.BookingConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start ProcessBookingQueue...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Printf("Error: ProcessBookingQueue failed: %s", err.Error())
			continue
		case msg := <-consumer.Messages():
			h.bookingUsecase.UpOffSet(ctx, msg.Offset+1)

			req := new(booking.BookingQueueMessage)
			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				continue
			}

			_, err := h.bookingUsecase.InsertBooking(ctx, req.FacilityName, &booking.CreateBookingRequest{
				UserId:          req.UserId,
				SlotId:          req.SlotId,
				BadmintonSlotId: req.BadmintonSlotId,
			})
			if err != nil {
				log.Printf("Error processing booking: %v", err)
				continue
			}

			log.Printf("ProcessBookingQueue | Topic(%s) | Offset(%d) Message(%s)\n", 
				msg.Topic, msg.Offset, string(msg.Value))

		case <-sigchan:
			log.Println("Stop ProcessBookingQueue...")
			return
		}
	}
}

func (h *bookingQueueHandler) HandleFacilityUpdates() {
	
	consumer, err := queue.ConnectConsumer(
		[]string{h.cfg.Kafka.Url},
		h.cfg.Kafka.ApiKey,
		h.cfg.Kafka.Secret,
	)
	if err != nil {
		log.Printf("Failed to connect to Kafka: %v", err)
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("facility.updated", 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to create partition consumer: %v", err)
		return
	}
	defer partitionConsumer.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from consumer: %v", err)
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received facility update: %s", string(msg.Value))
		case <-sigchan:
			return
		}
	}
}

func (h *bookingQueueHandler) PublishBookingCreated(booking *booking.Booking) error {
	message, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	return queue.PushMessageWithKeyToQueue(
		[]string{h.cfg.Kafka.Url},
		h.cfg.Kafka.ApiKey,
		h.cfg.Kafka.Secret,
		"booking.created",
		booking.UserId,
		message,
	)
}
