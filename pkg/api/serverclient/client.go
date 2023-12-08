package serverclient

import (
	"arch/internal/viewer/entity"
	"arch/proto/pb"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"time"
)

var emptyData []byte

const (
	natsRequestTimeout = 10 * time.Second
	emptyError         = ""
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

const getAllMessagesSubject = "server.messages.get.all"

func (c *Client) GetAllMessages() ([]entity.Message, error) {
	respond, err := c.conn.NatsConn().Request(getAllMessagesSubject, emptyData, natsRequestTimeout)
	if err != nil {
		return []entity.Message{}, fmt.Errorf("c.conn.NatsConn().Request: %w", err)
	}

	var protoMessages pb.GetAllResponse
	if err = proto.Unmarshal(respond.Data, &protoMessages); err != nil {
		return []entity.Message{}, fmt.Errorf("proto.Unmarshal: %w", err)
	}

	if protoMessages.Error != emptyError {
		return []entity.Message{}, errors.New(protoMessages.Error)
	}

	return lo.Map[*pb.Message, entity.Message](
		protoMessages.Messages,
		func(item *pb.Message, _ int) entity.Message {
			return entity.MessageFromProto(item)
		},
	), nil
}
