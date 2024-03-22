package database

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Scheme       string
	UserName     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

func buildConnectionURL(config *DBConfig) string {
	dbURL := &url.URL{
		Scheme: config.Scheme,
		User:   url.UserPassword(config.UserName, config.Password),
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.DatabaseName,
	}
	return dbURL.String()

}

func NewConnection(config *DBConfig) (*pgxpool.Pool, error) {
	fmt.Println(buildConnectionURL(config))
	dbpool, err := pgxpool.New(context.Background(), buildConnectionURL(config))
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
