package service

import (
	"context"
	"encoding/json"
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/repository"
	"main/pkg/queue"

	"github.com/IBM/sarama"
)

type BookingQueueService struct {
    cfg      *config.Config
    repo     repository.BookingRepositoryService
    consumer sarama.Consumer
}

func NewBookingQueueService(cfg *config.Config, repo repository.BookingRepositoryService) (*BookingQueueService, error) {
    consumer, err := queue.ConnectConsumer(
        []string{cfg.Kafka.Url},
        cfg.Kafka.ApiKey,
        cfg.Kafka.Secret,
    )
    if err != nil {
        return nil, err
    }

    return &BookingQueueService{
        cfg:      cfg,
        repo:     repo,
        consumer: consumer,
    }, nil
}

func (s *BookingQueueService) Start(ctx context.Context) {
    partitionConsumer, err := s.consumer.ConsumePartition("bookings", 0, sarama.OffsetNewest)
    if err != nil {
        log.Printf("Error creating partition consumer: %v", err)
        return
    }
    defer partitionConsumer.Close()

    for {
        select {
        case msg := <-partitionConsumer.Messages():
            var bookingMsg booking.BookingQueueMessage
            if err := json.Unmarshal(msg.Value, &bookingMsg); err != nil {
                log.Printf("Error unmarshalling message: %v", err)
                continue
            }

			if err := s.repo.ProcessBookingQueue(ctx, s.cfg, &bookingMsg); err != nil {
                log.Printf("Error processing booking: %v", err)
                continue
            }

        case err := <-partitionConsumer.Errors():
            log.Printf("Error from consumer: %v", err)

        case <-ctx.Done():
            return
        }
    }
} 