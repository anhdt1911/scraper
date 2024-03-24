package router

import (
	"encoding/gob"
	"net/http"

	"github.com/anhdt1911/scraper/internal/authenticator"
	"github.com/anhdt1911/scraper/internal/middleware"
	"github.com/anhdt1911/scraper/internal/server"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func New(s *server.Server, auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// Define path
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
	})

	router.GET("/", auth.Login)
	router.GET("/callback", auth.Callback)
	router.GET("/logout", auth.Logout)

	router.GET("/result/:keyID", middleware.IsAuthenticated, s.GetSearchResultByKeyword)
	router.GET("/results/:userID", middleware.IsAuthenticated, s.GetSearchResultsByUserID)

	router.POST("scrape", middleware.IsAuthenticated, s.ScrapeResult)
	router.POST("/batch-scrape", middleware.IsAuthenticated, s.BatchScrape)

	return router
}
