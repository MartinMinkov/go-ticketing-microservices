package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type Publisher struct {
	ns      *nats.Conn
	js      nats.JetStreamContext
	ctx     context.Context
	subject Subjects
}

func NewPublisher(ns *nats.Conn, subject Subjects, ctx context.Context) *Publisher {
	js, err := ns.JetStream()
	if err != nil {
		log.Fatal().Msgf("jetstream.New error: %v", err)
	}

	return &Publisher{
		ns:      ns,
		js:      js,
		ctx:     ctx,
		subject: subject,
	}
}

func (p *Publisher) Publish(data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println("Publishing to subject: ", string(p.subject))
	fmt.Println("Data: ", string(dataBytes))
	if _, err := p.js.Publish(string(p.subject), dataBytes); err != nil {
		return err
	}
	return nil
}
