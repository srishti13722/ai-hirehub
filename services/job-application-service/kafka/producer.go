package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

func InitKafkaWriter(){
	broker := os.Getenv("KAFKA_BROKER")

	Writer = &kafka.Writer{
		Addr: kafka.TCP(broker),
		Balancer: &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	log.Println("Kafka producer initialized..")
}

func Publish(topic string, value interface{}) error{
	msgBytes, err := json.Marshal(value)
	if err != nil{
		return err
	}

	return Writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   []byte(time.Now().Format(time.RFC3339)),
		Value: msgBytes,
	})
}