package auth

import (
	"arch/internal/server/entity"
	"arch/internal/server/service/auth/repo"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
)

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrBadPermission = errors.New("bad permission")
)

type authRepo interface {
	Credentials(ctx context.Context, credentials entity.UserCredentials) (entity.UserCredentials, error)
	CheckPermission(ctx context.Context, permission int, resource entity.ResourceInfo) error
}

type Service struct {
	authRepo authRepo
	logger   *zerolog.Logger
}

func New(repo authRepo, logger *zerolog.Logger) *Service {
	return &Service{
		authRepo: repo,
		logger:   logger,
	}
}

func (s *Service) Login(ctx context.Context, credentials entity.UserCredentials, resource entity.ResourceInfo) error {
	userCredentials, err := s.authRepo.Credentials(ctx, credentials)
	if err != nil {
		if errors.Is(err, repo.ErrNotAllowed) {
			return ErrUnauthorized
		}

		return fmt.Errorf("s.authRepo.Credentials: %w", err)
	}

	if err = s.authRepo.CheckPermission(ctx, userCredentials.Permission, resource); err != nil {
		if errors.Is(err, repo.ErrNotAllowed) {
			return ErrBadPermission
		}

		return fmt.Errorf("s.authRepo.CheckPermission: %w", err)
	}

	return nil
}
