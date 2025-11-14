package service

import (
	"tariff-api/internal/model"
	"tariff-api/internal/repository"

	"github.com/pkg/errors"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(username, password string) (*model.TokenResponse, error) {
	tokenResponse, err := s.repo.Login(username, password)
	if err != nil {
		return nil, errors.Wrap(err, "service failed to login")
	}
	return tokenResponse, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*model.TokenResponse, error) {
	tokenResponse, err := s.repo.RefreshToken(refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "service failed to refresh token")
	}
	return tokenResponse, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	err := s.repo.Logout(refreshToken)
	if err != nil {
		return errors.Wrap(err, "service failed to logout")
	}
	return nil
}
