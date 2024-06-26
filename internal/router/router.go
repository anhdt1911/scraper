package router

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/anhdt1911/scraper/internal/authenticator"
	"github.com/anhdt1911/scraper/internal/config"
	"github.com/anhdt1911/scraper/internal/middleware"
	"github.com/anhdt1911/scraper/internal/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func New(s *server.Server, auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.UIDomain, "https://" + config.AuthDomain},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
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
	router.GET("/user", middleware.IsAuthenticated, auth.GetUser)

	router.GET("/result/:keyID", middleware.IsAuthenticated, s.GetSearchResultByKeyword)
	router.GET("/results/:userID", middleware.IsAuthenticated, s.GetSearchResultsByUserID)

	router.POST("scrape", middleware.IsAuthenticated, s.ScrapeResult)
	router.POST("/batch-scrape", middleware.IsAuthenticated, s.BatchScrape)

	return router
}
