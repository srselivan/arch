package nats

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"time"
)

const (
	connectWait  = 10 * time.Second
	pingInterval = 10
	pingMaxOut   = 5

	stanClusterID = "nats-streaming"
	natsURL       = "nats://127.0.0.1:14222,nats://127.0.0.1:24222,nats://127.0.0.1:34222"
)

type Config struct {
	ClientID string
}

func New(cfg Config) (stan.Conn, error) {
	conn, err := stan.Connect(
		stanClusterID,
		cfg.ClientID,
		stan.NatsURL(natsURL),
		stan.ConnectWait(connectWait),
		stan.Pings(pingInterval, pingMaxOut),
	)
	if err != nil {
		return nil, fmt.Errorf("stan.Connect: %w", err)
	}
	return conn, nil
}
