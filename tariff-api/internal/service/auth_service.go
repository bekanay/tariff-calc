package service

import (
	"tariff-api/internal/model"
	"tariff-api/internal/repository"

	"github.com/pkg/errors"
)

type AuthService struct {
	repo         *repository.AuthRepository
	requiredRole string
}

func NewAuthService(repo *repository.AuthRepository, requiredRole string) *AuthService {
	return &AuthService{repo: repo, requiredRole: requiredRole}
}

var ErrMissingRole = errors.New("required role missing")

func (s *AuthService) Login(username, password string) (*model.TokenResponse, error) {
	tokenResponse, err := s.repo.Login(username, password)
	if err != nil {
		return nil, errors.Wrap(err, "service failed to login")
	}

	if s.requiredRole != "" {
		hasRole, err := s.HasRole(tokenResponse.AccessToken, s.requiredRole)
		if err != nil {
			return nil, errors.Wrap(err, "service failed to verify role")
		}
		if !hasRole {
			return nil, ErrMissingRole
		}
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

func (s *AuthService) HasRole(accessToken, role string) (bool, error) {
	hasRole, err := s.repo.HasRole(accessToken, role)
	if err != nil {
		return false, errors.Wrap(err, "service failed to check role")
	}
	return hasRole, nil
}
