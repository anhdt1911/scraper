package main

import (
	"log"

	"github.com/anhdt1911/scraper/internal/authenticator"
	"github.com/anhdt1911/scraper/internal/config"
	"github.com/anhdt1911/scraper/internal/database"
	"github.com/anhdt1911/scraper/internal/router"
	"github.com/anhdt1911/scraper/internal/scraper"
	"github.com/anhdt1911/scraper/internal/server"
)

func main() {
	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Fail to init authenticator: %v", err)
	}

	// Init scraper
	scrpr := scraper.New()

	// Init database connection
	db, err := database.NewConnection(&database.DBConfig{
		Scheme:       config.DBScheme,
		UserName:     config.DBUserName,
		Password:     config.DBPassword,
		Host:         config.DBHost,
		Port:         config.DBPort,
		DatabaseName: config.DBName,
	})
	if err != nil {
		panic("error initialize database connection")
	}
	defer db.Close()

	s := server.New(db, scrpr)

	router := router.New(s, auth)
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("An error occurs with the server: %v", err)
	}
}
