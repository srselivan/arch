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
	Login(ctx context.Context, credentials entity.UserCredentials) (int, error)
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
	perm, err := s.authRepo.Login(ctx, credentials)
	if err != nil {
		if errors.Is(err, repo.ErrNotAllowed) {
			return ErrUnauthorized
		}

		return fmt.Errorf("s.authRepo.Login: %w", err)
	}

	if err = s.authRepo.CheckPermission(ctx, perm, resource); err != nil {
		if errors.Is(err, repo.ErrNotAllowed) {
			return ErrBadPermission
		}

		return fmt.Errorf("s.authRepo.CheckPermission: %w", err)
	}

	return nil
}
