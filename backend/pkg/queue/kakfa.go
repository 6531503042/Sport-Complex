package queue

// import (
// 	"crypto/tls"
// 	"encoding/json"
// 	"errors"
// 	"log"

// 	"github.com/IBM/sarama"
// 	"github.com/go-playground/validator"
// )

// func ConnectProducer(brokerUrls []string, apiKey, secret string) (sarama.SyncProducer, error) {
// 	config := sarama.NewConfig()
// 	if apiKey != "" && secret != "" {
// 		config.Net.SASL.Enable = true
// 		config.Net.SASL.User = apiKey
// 		config.Net.SASL.Password = secret
// 		config.Net.SASL.Mechanism = "PLAIN"
// 		config.Net.SASL.Handshake = true
// 		config.Net.SASL.Version = sarama.SASLHandshakeV1
// 		config.Net.TLS.Enable = true
// 		config.Net.TLS.Config = &tls.Config{
// 			InsecureSkipVerify: true,
// 			ClientAuth: tls.NoClientCert,
// 		}
// 	}
// 	config.Producer.Return.Successes = true
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Retry.Max = 3

// 	producer, err := sarama.NewSyncProducer(brokerUrls, config)
// 	if err != nil {
// 		log.Printf("Error: Failed to connect to producer: %s", err.Error())
// 		return nil, errors.New("error: failed to connect to producer")
// 	}
// 	return producer, nil
// }

// func ConnectConsumer(brokerUrls []string, apiKey, secret string) (sarama.Consumer, error) {
// 	config := sarama.NewConfig()
// 	if apiKey != "" && secret != "" {
// 		config.Net.SASL.Enable = true
// 		config.Net.SASL.User = apiKey
// 		config.Net.SASL.Password = secret
// 		config.Net.SASL.Mechanism = "PLAIN"
// 		config.Net.SASL.Handshake = true
// 		config.Net.SASL.Version = sarama.SASLHandshakeV1
// 		config.Net.TLS.Enable = true
// 		config.Net.TLS.Config = &tls.Config{
// 			InsecureSkipVerify: true,
// 			ClientAuth: tls.NoClientCert,
// 		}
// 	}
// 	config.Consumer.Return.Errors = true
// 	config.Consumer.Fetch.Max = 3

// 	consumer, err := sarama.NewConsumer(brokerUrls, config)
// 	if err != nil {
// 		log.Printf("Error: Failed to connect to consumer: %s", err.Error())
// 		return nil, errors.New("error: failed to connect to consumer")
// 	}

// 	return consumer, nil
// }

// func PushMessageWithKeyToQueue(brokerUrls []string, apiKey, secret, topic, key string, message []byte) error {
// 	producer, err := ConnectProducer(brokerUrls, apiKey, secret)
// 	if err != nil {
// 		log.Printf("Error: Failed to connect to producer: %s", err.Error())
// 		return errors.New("error: failed to connect to producer")
// 	}
// 	defer producer.Close()

// 	msg := &sarama.ProducerMessage{
// 		Topic: topic,
// 		Value: sarama.StringEncoder(message),
// 		Key:   sarama.StringEncoder(key),
// 	}

// 	partition, offset, err := producer.SendMessage(msg)
// 	if err != nil {
// 		log.Printf("Error: Failed to push message to queue: %s", err.Error())
// 		return errors.New("error: failed to push message to queue")
// 	}

// 	log.Printf("Message pushed to queue: partition=%d, offset=%d", partition, offset)
// 	return nil
// }

// func DecodeMessage(obj any, value []byte) error {
// 	if err := json.Unmarshal(value, &obj); err != nil {
// 		log.Printf("Error: Failed to decode message: %s", err.Error())
// 		return errors.New("error: failed to decode message")
// 	}

// 	validate := validator.New()
// 	if err := validate.Struct(obj); err != nil {
// 		log.Printf("Error: Failed to validate message: %s", err.Error())
// 		return errors.New("error: failed to validate message")
// 	}

// 	return nil
// }