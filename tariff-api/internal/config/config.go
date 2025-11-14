package config

import (
	"log"
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
	}
	Dsn string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &Config{}

	cfg.KeyCloak.KeycloakURL = os.Getenv("KEYCLOAK_URL")
	cfg.KeyCloak.ClientID = os.Getenv("CLIENT_ID")
	cfg.KeyCloak.ClientSecret = os.Getenv("CLIENT_SECRET")
	cfg.KeyCloak.Realm = os.Getenv("REALM")
	cfg.KeyCloak.RedirectURI = os.Getenv("REDIRECT_URI")

	cfg.Dsn = os.Getenv("DB_DSN")

	return cfg
}
