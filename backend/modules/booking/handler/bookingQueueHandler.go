package handler

import (
	"context"
	"encoding/json"
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
		AddBooking()
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

func (h *bookingQueueHandler) AddBooking() {
    ctx := context.Background()

    consumer, err := h.BookingConsumer(ctx)
    if err != nil {
        log.Println("Error: AddBooking failed: ", err.Error())
        return
    }
    defer consumer.Close()

    log.Println("Start listening for new booking messages...")

    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    for {
        select {
        case err := <-consumer.Errors():
            log.Println("Error: AddBooking failed: ", err.Error())
            continue
        case msg := <-consumer.Messages():
            if string(msg.Key) == "booking" {
                // Update the offset
                if err := h.bookingUsecase.UpOffSet(ctx, msg.Offset+1); err != nil {
                    log.Println("Error: failed to update offset: ", err.Error())
                    continue
                }

                // Decode the booking request from Kafka message
                req := new(booking.Booking)
                if err := json.Unmarshal(msg.Value, req); err != nil {
                    log.Println("Error: failed to unmarshal booking message: ", err.Error())
                    continue
                }

                // Process the booking request (inserting it into the system)
                _, err := h.bookingUsecase.InsertBooking(ctx, req.UserId, req.SlotId)
                if err != nil {
                    log.Println("Error: failed to insert booking: ", err.Error())
                    continue
                }

                log.Println("Successfully processed and inserted booking for user:", req.UserId)
            }
        case sig := <-sigchan:
            log.Printf("Received signal %v, shutting down...", sig)
            return
        }
    }
}
