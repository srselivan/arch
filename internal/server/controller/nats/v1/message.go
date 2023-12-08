package v1

import (
	"arch/internal/server/entity"
	"arch/proto/pb"
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

const (
	getMessagesSubject = "server.messages.get.all"
	getMessagesGroup   = "server.messages.get.all_group"
)

func (h *Handler) InitMessagesRoutes() {
	if _, err := h.conn.NatsConn().QueueSubscribe(getMessagesSubject, getMessagesGroup, h.getAllMessages); err != nil {
		h.log.Fatal().Err(err)
	}
}

func (h *Handler) getAllMessages(msg *nats.Msg) {
	errorMessage := ""

	messages, err := h.messageService.GetAll(context.Background())
	if err != nil {
		h.log.Error().Err(err)
		errorMessage = err.Error()
	}

	protoMessages := lo.Map[entity.Message, *pb.Message](
		messages,
		func(item entity.Message, _ int) *pb.Message {
			return item.ToProto()
		},
	)

	responseBytes, err := proto.Marshal(&pb.GetAllResponse{
		Messages: protoMessages,
		Error:    errorMessage,
	})
	if err != nil {
		h.log.Error().Err(err)
	}

	if err = msg.Respond(responseBytes); err != nil {
		h.log.Error().Err(err)
	}
}
