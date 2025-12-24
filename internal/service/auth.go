package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kstsm/wb-warehouse-control/internal/models"
	"github.com/kstsm/wb-warehouse-control/pkg/jwt"
)

func (s *Service) SignInOrSignUp(ctx context.Context, userName, role string) (string, string, error) {
	user := models.User{
		ID:        uuid.New(),
		Name:      userName,
		Role:      role,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	resultUser, err := s.repo.GetOrCreateUser(ctx, user)
	if err != nil {
		return "", "", err
	}

	jwtRole := jwt.Role(resultUser.Role)
	token, err := s.tokenGenerator.GenerateToken(resultUser.ID, jwtRole)
	if err != nil {
		return "", "", err
	}

	return token, resultUser.Role, nil
}
