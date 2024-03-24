package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	DBScheme   = configVar("POSTGRES_SCHEMA")
	DBUserName = configVar("POSTGRES_USER")
	DBPassword = configVar("POSTGRES_PASSWORD")
	DBHost     = configVar("POSTGRES_HOST")
	DBPort     = configVar("POSTGRES_PORT")
	DBName     = configVar("POSTGRES_DB")
)

var (
	AuthDomain       = configVar("AUTH0_DOMAIN")
	AuthClientID     = configVar("AUTH0_CLIENT_ID")
	AuthClientSecret = configVar("AUTH0_CLIENT_SECRET")
	AuthCallbackURL  = configVar("AUTH0_CALLBACK_URL")
)

func configVar(name string) string {
	// Todo move load env to main
	err := godotenv.Load("./.env")
	if err != nil {
		panic("error loading env variables")
	}
	return os.Getenv(name)
}
