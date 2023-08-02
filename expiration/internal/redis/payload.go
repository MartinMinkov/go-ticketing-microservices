package redis

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const ExpirationType = "expiration"

type Payload struct {
	OrderId string `json:"order_id"`
}

func CreateExpirationTask(orderId string) *asynq.Task {
	payload, err := json.Marshal(Payload{OrderId: orderId})
	if err != nil {
		log.Fatal(err)
	}
	return asynq.NewTask(ExpirationType, payload)
}

func EnqueueTask(client *asynq.Client, task *asynq.Task, duration time.Duration) error {
	_, err := client.Enqueue(task, asynq.ProcessIn(duration))
	return err
}

func EnqueueCreateExpiration(client *asynq.Client, orderId string, duration time.Duration) error {
	task := CreateExpirationTask(orderId)
	return EnqueueTask(client, task, duration)
}

func ParsePayload(payload []byte) (string, error) {
	var p Payload
	err := json.Unmarshal(payload, &p)
	if err != nil {
		return "", err
	}
	return p.OrderId, nil
}
