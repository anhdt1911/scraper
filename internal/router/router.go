package router

import (
	"github.com/anhdt1911/scraper/internal/server"
	"github.com/gin-gonic/gin"
)

func New(s *server.Server) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	// Define path
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/result/:keyID", s.GetSearchResultByKeyword)
	router.GET("/results/:userID", s.GetSearchResultsByUserID)

	router.POST("scrape", s.ScrapeResult)
	router.POST("/batch-scrape", s.BatchScrape)

	return router
}
