package events

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

func ConnectWithRetry(address string, retryInterval time.Duration, maxDuration time.Duration) (*nats.Conn, error) {
	startTime := time.Now()
	for {
		elapsedTime := time.Since(startTime)
		if elapsedTime > maxDuration {
			return nil, nats.ErrNoServers
		}

		nc, err := nats.Connect(address)
		if err == nil {
			return nc, nil
		}

		log.Info().Msgf("Failed to connect to NATS at %s. Retrying...", address)
		time.Sleep(retryInterval)
	}
}
