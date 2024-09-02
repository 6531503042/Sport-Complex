package kafka

import (
	"context"
	"encoding/json"
	"log"
	"main/modules/booking/models"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
    writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
    return &KafkaProducer{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Topic:    topic,
            Balancer: &kafka.LeastBytes{},
        },
    }
}

func (p *KafkaProducer) PublishBookingEvent(ctx context.Context, booking models.Booking) error {
    msg, err := json.Marshal(booking)
    if err != nil {
        return err
    }

    err = p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(booking.ID.Hex()),
        Value: msg,
    })
    if err != nil {
        return err
    }

    log.Printf("Booking event published: %s", booking.ID.Hex())
    return nil
}
