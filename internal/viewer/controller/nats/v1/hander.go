package v1

import (
	"arch/internal/viewer/entity"
	"context"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

type messageService interface {
	Save(ctx context.Context, message entity.Message)
	Get(cxt context.Context, id string) (entity.Message, error)
	GetAll(ctx context.Context) ([]entity.Message, error)
}

type consoleUpdater interface {
	UpdateScreen() error
}

type Config struct {
	Conn           stan.Conn
	Log            *zerolog.Logger
	MessageService messageService
	ConsoleUpdater consoleUpdater
}

type Handler struct {
	conn           stan.Conn
	log            *zerolog.Logger
	messageService messageService
	consoleUpdater consoleUpdater
}

func New(cfg Config) *Handler {
	return &Handler{
		conn:           cfg.Conn,
		log:            cfg.Log,
		messageService: cfg.MessageService,
		consoleUpdater: cfg.ConsoleUpdater,
	}
}
