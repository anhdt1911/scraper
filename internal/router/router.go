package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	// Define path
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/upload", func(ctx *gin.Context) {
		file, _ := ctx.FormFile("file")
		log.Println(file.Filename)
		dir, _ := os.Getwd()
		fmt.Println(dir)
		ctx.SaveUploadedFile(file, dir+"/"+file.Filename)
		ctx.String(http.StatusOK, fmt.Sprintf("'%s uploaded!", file.Filename))
	})

	return router
}
