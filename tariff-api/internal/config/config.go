package config

import (
	"fmt"
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
	PgDsn    string
	Db2Dsn   string
	LogLevel string
	GinMode  string
	Port     string
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

	cfg.PgDsn = os.Getenv("PG_DSN")
	cfg.Db2Dsn = os.Getenv("DB2_DSN")
	cfg.LogLevel = os.Getenv("LOG_LEVEL")
	cfg.GinMode = os.Getenv("GIN_MODE")
	cfg.Port = os.Getenv("PORT")

	if cfg.Port == "" {
		cfg.Port = "8081"
	}

	if cfg.PgDsn == "" {
		return nil, fmt.Errorf("PG_DSN is required")
	}
	if cfg.Db2Dsn == "" {
		return nil, fmt.Errorf("DB2_DSN is required")
	}

	return cfg, nil
}
