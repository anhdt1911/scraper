package router

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	router := gin.Default()

	// Define path
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
