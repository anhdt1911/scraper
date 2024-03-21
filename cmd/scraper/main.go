package main

import (
	"log"

	"github.com/anhdt1911/scraper/internal/router"
)

func main() {
	router := router.New()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("An error occurs with the server: %v", err)
	}
}
