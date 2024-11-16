package queue

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func ConnectProducer(brokerUrls []string, apiKey, secret string) (sarama.SyncProducer, error) {
    config := sarama.NewConfig()
    config.Version = sarama.V2_1_0_0
    
    // Set required configurations
    config.Producer.Return.Successes = true
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Retry.Backoff = 500 * time.Millisecond
    
    // Set up SASL if credentials are provided
    if apiKey != "" && secret != "" {
        config.Net.SASL.Enable = true
        config.Net.SASL.User = apiKey
        config.Net.SASL.Password = secret
        config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
        config.Net.SASL.Handshake = true
        config.Net.TLS.Enable = true
        config.Net.TLS.Config = &tls.Config{
            InsecureSkipVerify: true,
        }
    }

    // Create producer with retry logic
    var producer sarama.SyncProducer
    var err error
    
    for i := 0; i < 3; i++ {
        producer, err = sarama.NewSyncProducer(brokerUrls, config)
        if err == nil {
            return producer, nil
        }
        log.Printf("Failed to connect to Kafka (attempt %d): %v", i+1, err)
        time.Sleep(time.Second * 2)
    }

    return nil, err
}

func PushMessageWithKeyToQueue(brokerUrls []string, apiKey, secret, topic, key string, message []byte) error {
    config := sarama.NewConfig()
    config.Version = sarama.V2_1_0_0
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Return.Successes = true
    config.Producer.Retry.Max = 5

    producer, err := sarama.NewSyncProducer(brokerUrls, config)
    if err != nil {
        log.Printf("Failed to create producer: %v", err)
        return fmt.Errorf("error: failed to connect to producer: %w", err)
    }
    defer producer.Close()

    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(key),
        Value: sarama.ByteEncoder(message),
    }

    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        return fmt.Errorf("error: failed to send message: %w", err)
    }

    log.Printf("Message sent successfully to partition %d at offset %d", partition, offset)
    return nil
}

func ConnectConsumer(brokerUrls []string, apiKey, secret string) (sarama.Consumer, error) {
    config := sarama.NewConfig()
    config.Version = sarama.V2_1_0_0
    config.Consumer.Return.Errors = true
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    
    // Add debug logging
    config.ClientID = "debug-client"
    
    consumer, err := sarama.NewConsumer(brokerUrls, config)
    if err != nil {
        log.Printf("Failed to create consumer: %v", err)
        return nil, err
    }
    
    log.Printf("Successfully connected to Kafka")
    return consumer, nil
}

func DecodeMessage(obj interface{}, value []byte) error {
    if err := json.Unmarshal(value, obj); err != nil {
        log.Printf("Failed to decode message: %v", err)
        return err
    }
    return nil
} 