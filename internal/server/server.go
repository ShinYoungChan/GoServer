package server

import (
	"gin/internal/routes"

	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	routes.Register(r)

	return r
}
