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

func configVar(name string) string {
	err := godotenv.Load("./.env")
	if err != nil {
		panic("error loading env variables")
	}
	return os.Getenv(name)
}
