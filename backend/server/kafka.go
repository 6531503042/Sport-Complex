package server

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

type KafkaServerService interface {
	Start() error
	Stop()
	RegisterHandler(topic string, handler func(data []byte) error)
}

func NewKafkaServer(cfg *config.Config) KafkaServerService {
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
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0  // Set to a lower version
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{s.cfg.Kafka.Url}, config)
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