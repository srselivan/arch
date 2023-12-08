package viewerclient

import (
	"arch/internal/server/entity"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

type Client struct {
	log  *zerolog.Logger
	conn stan.Conn
}

func New(log *zerolog.Logger, conn stan.Conn) *Client {
	return &Client{
		log:  log,
		conn: conn,
	}
}

const publishMessageSubject = "viewer.updates.message"

func (c *Client) PublishMessage(message entity.Message) error {
	c.log.Debug().
		Str("subject", publishMessageSubject).
		Any("message", message).
		Msg("Publish message to nats")

	bytes, err := proto.Marshal(message.ToProto())
	if err != nil {
		return fmt.Errorf("prtot.Marshal: %w", err)
	}

	if err = c.conn.NatsConn().Publish(publishMessageSubject, bytes); err != nil {
		return fmt.Errorf("c.conn.Publish: %w", err)
	}

	return nil
}
