package model

import (
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

type KeycloakAccessTokenClaims struct {
	jwt.RegisteredClaims
	RealmAccess    RealmAccess               `json:"realm_access"`
	ResourceAccess map[string]ResourceAccess `json:"resource_access"`
}

func (claims KeycloakAccessTokenClaims) HasRole(role string) bool {
	if role == "" {
		return false
	}

	if containsRole(claims.RealmAccess.Roles, role) {
		return true
	}

	for _, access := range claims.ResourceAccess {
		if containsRole(access.Roles, role) {
			return true
		}
	}

	return false
}

func containsRole(roles []string, needle string) bool {
	for _, role := range roles {
		if strings.EqualFold(role, needle) {
			return true
		}
	}
	return false
}
