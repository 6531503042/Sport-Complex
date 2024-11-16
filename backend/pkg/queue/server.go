package queue

import (
	"log"
	"main/config"
	"sync"

	"github.com/IBM/sarama"
)

type KafkaServer struct {
	cfg      *config.Config
	consumer sarama.Consumer
	topics   []string
	handlers map[string]func(data []byte) error
	mu       sync.RWMutex
}

func NewKafkaServer(cfg *config.Config) *KafkaServer {
	return &KafkaServer{
		cfg:      cfg,
		topics:   make([]string, 0),
		handlers: make(map[string]func(data []byte) error),
	}
}

func (s *KafkaServer) RegisterHandler(topic string, handler func(data []byte) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.topics = append(s.topics, topic)
	s.handlers[topic] = handler
}

func (s *KafkaServer) Start() error {
	consumer, err := ConnectConsumer(
		[]string{s.cfg.Kafka.Url},
		s.cfg.Kafka.ApiKey,
		s.cfg.Kafka.Secret,
	)
	if err != nil {
		return err
	}
	s.consumer = consumer

	for _, topic := range s.topics {
		go s.consumeTopic(topic)
	}

	return nil
}

func (s *KafkaServer) Stop() {
	if s.consumer != nil {
		if err := s.consumer.Close(); err != nil {
			log.Printf("Error closing Kafka consumer: %v", err)
		}
	}
}

func (s *KafkaServer) consumeTopic(topic string) {
	partitionConsumer, err := s.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Error creating partition consumer for topic %s: %v", topic, err)
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			s.mu.RLock()
			handler := s.handlers[topic]
			s.mu.RUnlock()

			if handler != nil {
				if err := handler(msg.Value); err != nil {
					log.Printf("Error handling message for topic %s: %v", topic, err)
				}
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from consumer for topic %s: %v", topic, err)
		}
	}
} 