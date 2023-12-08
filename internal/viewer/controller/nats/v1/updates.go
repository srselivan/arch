package v1

import (
	"arch/internal/viewer/entity"
	"arch/proto/pb"
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

const (
	updatesSubject = "viewer.updates.message"
)

func (h *Handler) InitUpdatesRoutes() {
	if _, err := h.conn.NatsConn().Subscribe(updatesSubject, h.messageUpdate); err != nil {
		h.log.Fatal().Err(err)
	}
}

func (h *Handler) messageUpdate(msg *nats.Msg) {
	var protoMessage pb.Message
	if err := proto.Unmarshal(msg.Data, &protoMessage); err != nil {
		h.log.Error().Err(err).Send()
		return
	}

	h.messageService.Save(context.Background(), entity.MessageFromProto(&protoMessage))

	if err := h.consoleUpdater.UpdateScreen(); err != nil {
		h.log.Error().Err(err).Send()
		return
	}
}
