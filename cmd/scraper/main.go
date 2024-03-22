package main

import (
	"log"

	"github.com/anhdt1911/scraper/internal/config"
	"github.com/anhdt1911/scraper/internal/database"
	"github.com/anhdt1911/scraper/internal/router"
)

func main() {
	_, err := database.NewConnection(&database.DBConfig{
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
	router := router.New()
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("An error occurs with the server: %v", err)
	}
}
