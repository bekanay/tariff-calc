package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	KeyCloak struct {
		KeycloakURL  string
		ClientID     string
		ClientSecret string
		Realm        string
		RedirectURI  string
		RequiredRole string
	}
	Dsn      string
	LogLevel string
	GinMode  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}

	cfg.KeyCloak.KeycloakURL = os.Getenv("KEYCLOAK_URL")
	cfg.KeyCloak.ClientID = os.Getenv("CLIENT_ID")
	cfg.KeyCloak.ClientSecret = os.Getenv("CLIENT_SECRET")
	cfg.KeyCloak.Realm = os.Getenv("REALM")
	cfg.KeyCloak.RedirectURI = os.Getenv("REDIRECT_URI")
	cfg.KeyCloak.RequiredRole = os.Getenv("REQUIRED_ROLE")

	cfg.Dsn = os.Getenv("DB_DSN")
	cfg.LogLevel = os.Getenv("LOG_LEVEL")
	cfg.GinMode = os.Getenv("GIN_MODE")

	return cfg, nil
}
