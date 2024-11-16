package handlers

import (
	"encoding/json"
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/user/usecase"
	"main/pkg/queue"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type (
	UserQueueHandlerService interface {
		ProcessBookingEvents()
		PublishUserUpdated(userId string, status string) error
	}

	userQueueHandler struct {
		cfg          *config.Config
		userUsecase usecase.UserUsecaseService
	}
)

func NewUserQueueHandler(cfg *config.Config, userUsecase usecase.UserUsecaseService) UserQueueHandlerService {
	return &userQueueHandler{
		cfg:          cfg,
		userUsecase: userUsecase,
	}
}

func (h *userQueueHandler) ProcessBookingEvents() {
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
			log.Printf("Processing booking event for user: %+v", bookingEvent)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from consumer: %v", err)
		case <-sigchan:
			return
		}
	}
}

func (h *userQueueHandler) PublishUserUpdated(userId string, status string) error {
	message, err := json.Marshal(map[string]string{
		"user_id": userId,
		"status":  status,
	})
	if err != nil {
		return err
	}

	return queue.PushMessageWithKeyToQueue(
		[]string{h.cfg.Kafka.Url},
		h.cfg.Kafka.ApiKey,
		h.cfg.Kafka.Secret,
		"user.updated",
		userId,
		message,
	)
}
