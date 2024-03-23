package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anhdt1911/scraper/internal/config"
	"github.com/anhdt1911/scraper/internal/database"
	"github.com/anhdt1911/scraper/internal/scraper"
	"github.com/anhdt1911/scraper/internal/server"
)

func main() {
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

	_ = server.New(db, scrpr)

	// router := router.New()
	// if err := router.Run(":3000"); err != nil {
	// 	log.Fatalf("An error occurs with the server: %v", err)
	// }
	var res scraper.SearchResult
	err = db.QueryRow(context.Background(), "SELECT * FROM search_result WHERE id = 1").
		Scan(&res.ID, &res.TotalSearchResult, &res.HtmlContent, &res.Keyword)
	if err != nil {
		fmt.Println(err)
	}
	resD, _ := json.Marshal(res)
	fmt.Println(string(resD))
}
