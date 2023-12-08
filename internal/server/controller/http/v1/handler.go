package v1

import (
	"arch/internal/server/entity"
	"context"
	"github.com/rs/zerolog"
	"text/template"
)

type messageService interface {
	Save(ctx context.Context, message entity.Message)
	Get(cxt context.Context, id string) (entity.Message, error)
	GetAll(ctx context.Context) ([]entity.Message, error)
}

type Config struct {
	Log            *zerolog.Logger
	MessageService messageService
}

type Handler struct {
	log            *zerolog.Logger
	template       *template.Template
	messageService messageService
}

func New(cfg Config) *Handler {
	return &Handler{
		log:            cfg.Log,
		template:       template.Must(template.ParseFiles("internal/server/controller/http/template/index.html")),
		messageService: cfg.MessageService,
	}
}
