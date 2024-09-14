package handler

import (
	"context"
	"log"
	"main/config"
	"main/modules/booking/usecase"
	"main/pkg/queue"

	"github.com/IBM/sarama"
)

type (
	BookingQueueHttpHandlerService interface {

	}

	bookingQueueHandler struct {
		cfg              *config.Config
		bookingUsecase usecase.BookingUsecaseService
	}
)

func NewBookingQueueHandler(cfg *config.Config, bookingUsecase usecase.BookingUsecaseService) BookingQueueHttpHandlerService {
	return &bookingQueueHandler{cfg: cfg, bookingUsecase: bookingUsecase}
}

func (h * bookingQueueHandler) BookingConsumer (pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.bookingUsecase.GetOffSet(pctx)
	if err != nil {
		return nil, err
	}
	consumer, err := worker.ConsumePartition("booking", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("booking", 0, 0)
		if err != nil {
			log.Println("Error: BookingConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}