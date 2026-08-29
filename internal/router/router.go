package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine *gin.Engine
}

func NewApp(engine *gin.Engine) *App {
	return &App{
		Engine: engine,
	}
}

func (a *App) SetUpV1() *gin.RouterGroup {
	v1 := a.Engine.Group("/api/v1")

	v1.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    "pong",
			"error":   nil,
		})
	})

	return v1
}
