package events

import (
	"bytes"
	"context"
	"encoding/gob"

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
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}
	if _, err := p.js.Publish(string(p.subject), buf.Bytes()); err != nil {
		return err
	}
	return nil
}
