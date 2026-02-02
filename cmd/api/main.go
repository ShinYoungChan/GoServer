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
	// db 연결 코드 직후 (예: db, err := gorm.Open(...))
	models.DB.AutoMigrate(&models.Article{}, &models.Comment{})

	fmt.Println("Route Start")
	srv := server.NewEngine()
	srv.Run()
}
