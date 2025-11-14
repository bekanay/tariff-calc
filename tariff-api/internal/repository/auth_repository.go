package repository

import (
	"context"
	"tariff-api/internal/config"
	"tariff-api/internal/model"

	"github.com/Nerzal/gocloak/v11"
	"github.com/pkg/errors"
)

type AuthRepository struct {
	client gocloak.GoCloak
	cfg    *config.Config
}

func NewAuthRepository(cfg *config.Config) *AuthRepository {
	client := gocloak.NewClient(cfg.KeyCloak.KeycloakURL)
	return &AuthRepository{
		client: client,
		cfg:    cfg,
	}
}

func (repo *AuthRepository) Login(username, password string) (*model.TokenResponse, error) {
	token, err := repo.client.Login(context.Background(), repo.cfg.KeyCloak.ClientID, repo.cfg.KeyCloak.ClientSecret, repo.cfg.KeyCloak.Realm, username, password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login")
	}
	return &model.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IdToken:      token.IDToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

func (repo *AuthRepository) RefreshToken(refreshToken string) (*model.TokenResponse, error) {
	token, err := repo.client.RefreshToken(context.Background(), refreshToken, repo.cfg.KeyCloak.ClientID, repo.cfg.KeyCloak.ClientSecret, repo.cfg.KeyCloak.Realm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to refresh token")
	}
	return &model.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IdToken:      token.IDToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

func (repo *AuthRepository) Logout(refreshToken string) error {
	err := repo.client.Logout(context.Background(), repo.cfg.KeyCloak.ClientID, repo.cfg.KeyCloak.ClientSecret, repo.cfg.KeyCloak.Realm, refreshToken)
	if err != nil {
		return errors.Wrap(err, "failed to logout")
	}
	return nil
}
