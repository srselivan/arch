package nats

import (
	v1 "arch/internal/server/controller/nats/v1"
	"arch/internal/server/entity"
	"context"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

type messageService interface {
	Save(ctx context.Context, message entity.Message)
	Get(cxt context.Context, id string) (entity.Message, error)
	GetAll(ctx context.Context) ([]entity.Message, error)
}

type Config struct {
	Conn           stan.Conn
	Log            *zerolog.Logger
	MessageService messageService
}

type Server struct {
	conn           stan.Conn
	log            *zerolog.Logger
	messageService messageService
}

func New(cfg Config) *Server {
	return &Server{
		conn:           cfg.Conn,
		log:            cfg.Log,
		messageService: cfg.MessageService,
	}
}

func (s *Server) Run() {
	v1Handler := v1.New(v1.Config{
		Conn:           s.conn,
		Log:            s.log,
		MessageService: s.messageService,
	})

	v1Handler.InitMessagesRoutes()
}
