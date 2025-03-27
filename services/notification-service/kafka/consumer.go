package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/srishti13722/ai-hirehub/notification-service/ws"
)

func consume(topic string, handler func([]byte)) {
	broker := os.Getenv("KAFKA_BROKER")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: topic + "-consumer",
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}
		handler(m.Value)
	}
}

// Listen to application.created
func ConsumeApplicationCreated() {
	consume("job.application.created", func(value []byte) {
		var event map[string]interface{}
		json.Unmarshal(value, &event)
		recruiterID := event["recruiter_id"].(string)

		ws.SendToUser(recruiterID, value)
	})
}

// Listen to application.status.updated
func ConsumeApplicationStatusUpdated() {
	consume("job.application.status.updated", func(value []byte) {
		var event map[string]interface{}
		json.Unmarshal(value, &event)
		jobSeekerID := event["jobseeker_id"].(string)

		ws.SendToUser(jobSeekerID, value)
	})
}
