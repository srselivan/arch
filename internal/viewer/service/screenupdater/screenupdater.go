package screenupdater

import (
	"arch/internal/viewer/entity"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type messageRepo interface {
	Set(ctx context.Context, message entity.Message)
	Get(cxt context.Context, id string) (entity.Message, error)
	GetAll(ctx context.Context) ([]entity.Message, error)
}

type serverClient interface {
	GetAllMessages() ([]entity.Message, error)
}

type presenter interface {
	UpdateView(messages []entity.Message) error
}

type Service struct {
	log          *zerolog.Logger
	messageRepo  messageRepo
	serverClient serverClient
	presenter    presenter
}

func New(log *zerolog.Logger, repo messageRepo, client serverClient, presenter presenter) *Service {
	s := &Service{
		log:          log,
		messageRepo:  repo,
		serverClient: client,
		presenter:    presenter,
	}

	messages, err := client.GetAllMessages()
	if err != nil {
		log.Fatal().Err(err).Msg("An error occurred when get all messages from server")
	}

	for _, message := range messages {
		repo.Set(context.Background(), message)
	}

	if err = s.presenter.UpdateView(messages); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when update view")
	}

	return s
}

func (s *Service) UpdateScreen() error {
	messages, err := s.messageRepo.GetAll(context.Background())
	if err != nil {
		return fmt.Errorf("s.messageService.GetAll: %w", err)
	}

	if err = s.presenter.UpdateView(messages); err != nil {
		return fmt.Errorf("s.presenter.UpdateView: %w", err)
	}

	return nil
}
