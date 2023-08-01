package events

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

type Data struct {
	Id    string
	Title string
	Price int
}

type Event interface {
	Subject() Subjects
	Data() interface{}
}

type MessageParser interface {
	ParseMessage(msg jetstream.Msg) (interface{}, error)
}

type MessageHandler interface {
	OnMessage(data interface{}, msg jetstream.Msg) error
}

type Listener struct {
	Ns             *nats.Conn
	Js             jetstream.JetStream
	Ctx            context.Context
	subject        Subjects
	queueGroupName string
	ackWait        time.Duration
	parser         MessageParser
	handler        MessageHandler
}

func NewListener(ns *nats.Conn, subject Subjects, queueGroupName string, ackWait time.Duration, parser MessageParser, handler MessageHandler, ctx context.Context) *Listener {
	js, err := jetstream.New(ns)
	if err != nil {
		log.Fatal().Msgf("jetstream.New error: %v", err)
	}

	return &Listener{
		Ns:             ns,
		Js:             js,
		Ctx:            ctx,
		subject:        subject,
		queueGroupName: queueGroupName,
		ackWait:        ackWait,
		parser:         parser,
		handler:        handler,
	}
}

func (l *Listener) Listen() error {
	s, err := createStream(l.Ctx, l.Js, string(l.subject), l.queueGroupName)
	if err != nil {
		log.Info().Msgf("createStream error: %v", err)
	}

	cons, err := createConsumer(l.Ctx, s, l.queueGroupName, l.ackWait)
	if err != nil {
		log.Info().Msgf("createConsumer error: %v", err)
		return err
	}

	l.consume(cons, l.Ctx)

	<-l.Ctx.Done()
	return nil
}

func createStream(ctx context.Context, js jetstream.JetStream, subject string, queueGroupName string) (jetstream.Stream, error) {
	streamName := fmt.Sprintf("%s-%s", queueGroupName, subject)

	existingStream, err := js.Stream(ctx, streamName)
	if err != nil && existingStream != nil {
		log.Info().Msgf("returning existing stream: %s", streamName)
		return existingStream, nil
	}

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createConsumer(ctx context.Context, s jetstream.Stream, queueGroupName string, ackWait time.Duration) (jetstream.Consumer, error) {
	log.Info().Msgf("creating consumer: %s", queueGroupName)
	info, err := s.Info(ctx)
	if err != nil {
		log.Info().Msgf("s.Info error: %v", err)
		return nil, err
	}
	log.Info().Msgf("info: %v", info)
	cons, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       queueGroupName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       ackWait,
		MaxAckPending: 60,
	})
	if err != nil {
		return nil, err
	}
	return cons, nil
}

func (l *Listener) consume(j jetstream.Consumer, ctx context.Context) {
	cc, err := j.Consume(func(msg jetstream.Msg) {
		data, err := l.parser.ParseMessage(msg)
		if err != nil {
			log.Info().Msgf("parseMessage error: %v", err)
			return
		}
		err = l.handler.OnMessage(data, msg)
		if err != nil {
			log.Info().Msgf("onMessage error: %v", err)
			return
		}
	}, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		fmt.Println("consume error handler: ", err)
	}))

	if err != nil {
		log.Fatal().Msgf("consume error: %v", err)
	}

	<-ctx.Done()
	defer cc.Stop()
}
