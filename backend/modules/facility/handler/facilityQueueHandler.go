package handler

import (
	"encoding/json"
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/facility/usecase"
	"main/pkg/queue"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type (
	FacilityQueueHandlerService interface {
		ProcessBookingEvents()
		PublishFacilityUpdated(facilityId string, status string) error
	}

	facilityQueueHandler struct {
		cfg             *config.Config
		facilityUsecase usecase.FacilityUsecaseService
	}
)

func NewFacilityQueueHandler(cfg *config.Config, facilityUsecase usecase.FacilityUsecaseService) FacilityQueueHandlerService {
	return &facilityQueueHandler{
		cfg:             cfg,
		facilityUsecase: facilityUsecase,
	}
}

func (h *facilityQueueHandler) ProcessBookingEvents() {
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

	partitionConsumer, err := consumer.ConsumePartition("booking.created", 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to create partition consumer: %v", err)
		return
	}
	defer partitionConsumer.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var bookingEvent booking.Booking
			if err := json.Unmarshal(msg.Value, &bookingEvent); err != nil {
				log.Printf("Error unmarshalling booking event: %v", err)
				continue
			}
			log.Printf("Processing booking event: %+v", bookingEvent)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from consumer: %v", err)
		case <-sigchan:
			return
		}
	}
}

func (h *facilityQueueHandler) PublishFacilityUpdated(facilityId string, status string) error {
	message, err := json.Marshal(map[string]string{
		"facility_id": facilityId,
		"status":     status,
	})
	if err != nil {
		return err
	}

	return queue.PushMessageWithKeyToQueue(
		[]string{h.cfg.Kafka.Url},
		h.cfg.Kafka.ApiKey,
		h.cfg.Kafka.Secret,
		"facility.updated",
		facilityId,
		message,
	)
}
