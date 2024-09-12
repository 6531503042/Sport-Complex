package handler

import (
	"main/config"
	"main/modules/booking/usecase"
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

// func (h *bookingQueueHandler) BookingConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
// 	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
// 	if err != nil {
// 		log.Printf()
// 	}
// }