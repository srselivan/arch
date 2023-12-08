package message

import (
	"arch/internal/server/entity"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type messageRepo interface {
	Set(ctx context.Context, message entity.Message)
	Get(cxt context.Context, id string) (entity.Message, error)
	GetAll(ctx context.Context) ([]entity.Message, error)
}

type notifier interface {
	PublishMessage(message entity.Message) error
}

type Service struct {
	messageRepo messageRepo
	log         *zerolog.Logger
	notifier    notifier
}

func New(log *zerolog.Logger, repo messageRepo, notifier notifier) *Service {
	return &Service{
		messageRepo: repo,
		log:         log,
		notifier:    notifier,
	}
}

func (s *Service) Save(ctx context.Context, message entity.Message) {
	s.messageRepo.Set(ctx, message)

	go func() {
		if err := s.notifier.PublishMessage(message); err != nil {
			s.log.Error().Err(err).Send()
		}
	}()
}

func (s *Service) Get(ctx context.Context, id string) (entity.Message, error) {
	message, err := s.messageRepo.Get(ctx, id)
	if err != nil {
		return entity.Message{}, fmt.Errorf("s.messageRepo.Get: %w", err)
	}
	return message, nil
}

func (s *Service) GetAll(ctx context.Context) ([]entity.Message, error) {
	messages, err := s.messageRepo.GetAll(ctx)
	if err != nil {
		return []entity.Message{}, fmt.Errorf("s.messageRepo.GetAll: %w", err)
	}
	return messages, nil
}
