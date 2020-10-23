package producer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"strconv"
	log "strimzi-producer/internal/platform/logger"
	"time"
)

type client struct {
	sleepTime  time.Duration
	bufferSize uint
	writer     *kafka.Writer
}

func NewProducerClient(sleepTime time.Duration, bufferSize uint, writer *kafka.Writer) *client {
	return &client{
		sleepTime:  sleepTime,
		bufferSize: bufferSize,
		writer:     writer,
	}
}

func (c *client) Produce() {
	log.Logger.Info("Starting the Producer...")
	for {
		buffer := make([]kafka.Message, c.bufferSize)
		var j uint
		for j = 0; j < c.bufferSize; j++ {
			id := strconv.Itoa(rand.Int())
			msg := MyMessage{
				Id:      id,
				Message: fmt.Sprintf("Just a message with a random number: %d", rand.Int()),
			}
			bytes, err := proto.Marshal(&msg)
			if err != nil {
				log.Logger.Errorf("Error in marshalling the message: %s", err.Error())
				return
			}
			buffer[j] = kafka.Message{
				Key:   []byte(id),
				Value: bytes,
			}
		}
		err := c.writer.WriteMessages(context.Background(), buffer...)
		log.Logger.Debugf("Sent %d messages to the topic", len(buffer))
		if err != nil {
			log.Logger.Errorf("Error in producing message to Kafka: %s", err.Error())
			return
		}
		time.Sleep(c.sleepTime)
	}
}
