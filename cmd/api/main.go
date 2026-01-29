package main

import (
	"fmt"

	"gin/internal/models"
	"gin/internal/server"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	fmt.Println("DB Start")
	models.InitDB()

	fmt.Println("Route Start")
	srv := server.NewEngine()
	srv.Run()
}
