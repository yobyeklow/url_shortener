package routes

import (
	"url_shortener/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Routes) {
	r.Use(middleware.LoggerMiddleware(), middleware.ApiKeyMiddleware(), middleware.AuthMiddleWare(), middleware.RateLimiterMiddleware())
	api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(api)
	}
}
