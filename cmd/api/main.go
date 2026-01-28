package main

import (
	"fmt"

	"gin/internal/server"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	fmt.Println("Start")
	srv := server.NewEngine()
	srv.Run()
}
