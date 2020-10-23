package kafka

import (
	"github.com/segmentio/kafka-go"
	log "strimzi-producer/internal/platform/logger"
)

func Consumer(servers []string, groupId string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  servers, //[]string{"localhost:9092"},
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		Logger:      kafka.LoggerFunc(log.Logger.Infof),
		ErrorLogger: kafka.LoggerFunc(log.Logger.Errorf),
	})
}

func Producer(servers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(servers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireOne,
		Logger:      kafka.LoggerFunc(log.Logger.Infof),
		ErrorLogger: kafka.LoggerFunc(log.Logger.Errorf),
	}
}
