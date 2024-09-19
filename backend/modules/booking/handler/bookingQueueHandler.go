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

	"github.com/IBM/sarama"
)

type (
	BookingQueueHttpHandlerService interface {
		AddBooking(newBooking *booking.Booking) error
	}

	bookingQueueHandler struct {
		cfg              *config.Config
		bookingUsecase   usecase.BookingUsecaseService
	}
)

func NewBookingQueueHandler(cfg *config.Config, bookingUsecase usecase.BookingUsecaseService) BookingQueueHttpHandlerService {
	return &bookingQueueHandler{cfg: cfg, bookingUsecase: bookingUsecase}
}

// BookingConsumer listens to the Kafka "booking" topic, consumes messages,
// and processes bookings by updating the Kafka offset after each successful booking.
func (h *bookingQueueHandler) BookingConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	// Connect to Kafka consumer
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	// Get the last processed Kafka offset
	offset, err := h.bookingUsecase.GetOffSet(pctx)
	if err != nil {
		return nil, err
	}

	// Start consuming messages from the Kafka topic "booking"
	consumer, err := worker.ConsumePartition("booking", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("booking", 0, 0)
		if err != nil {
			log.Printf("Error consuming Kafka partition: %s", err.Error())
			return nil, err
		}
	}

	// Listen for system interrupts to stop the consumer gracefully
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer consumer.Close()

		for {
			select {
			case msg := <-consumer.Messages():
				log.Printf("Consumed message offset %d\n", msg.Offset)
				h.processBookingMessage(pctx, msg)

				// Update the offset after processing the message
				if err := h.bookingUsecase.UpOffSet(pctx, msg.Offset+1); err != nil {
					log.Printf("Error updating offset: %s", err.Error())
				}
			case <-signals:
				log.Println("Received interrupt signal, shutting down consumer.")
				return
			case err := <-consumer.Errors():
				log.Printf("Error consuming message: %s", err.Error())
			}
		}
	}()

	return consumer, nil
}

// processBookingMessage processes a single booking message.
func (h *bookingQueueHandler) processBookingMessage(ctx context.Context, msg *sarama.ConsumerMessage) {
	var bookingRequest booking.Booking
	if err := json.Unmarshal(msg.Value, &bookingRequest); err != nil {
		log.Printf("Error unmarshalling booking message: %s", err.Error())
		return
	}

	// Process the booking (you can modify this as per your logic)
	createdBooking, err := h.bookingUsecase.InsertBooking(ctx, bookingRequest.UserId, bookingRequest.SlotId)
	if err != nil {
		log.Printf("Error processing booking: %s", err.Error())
		return
	}

	log.Printf("Successfully processed booking for UserID %s, SlotID %s, BookingID %s", createdBooking.UserId, createdBooking.SlotId, createdBooking.Id.Hex())
}

func (h *bookingQueueHandler) AddBooking(newBooking *booking.Booking) error {
	// Connect to Kafka producer
	producer, err := queue.ConnectProducer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return fmt.Errorf("Error connecting to Kafka producer: %s", err.Error())
	}

	// Serialize the booking object to JSON
	bookingData, err := json.Marshal(newBooking)
	if err != nil {
		return fmt.Errorf("Error marshalling booking: %s", err.Error())
	}

	// Push the booking message to the Kafka topic "booking"
	msg := &sarama.ProducerMessage{
		Topic: "booking",
		Value: sarama.StringEncoder(bookingData),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Error sending booking to Kafka: %s", err.Error())
	}

	log.Println("Successfully added booking to Kafka")
	return nil
}
