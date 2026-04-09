package routes

import (
	"url_shortener/internal/handler"
	"url_shortener/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UrlRoutes struct {
	handler *handler.UrlsHandler
}

func NewUrlRoutes(handler *handler.UrlsHandler) *UrlRoutes {
	return &UrlRoutes{
		handler: handler,
	}
}

func (urlRoute *UrlRoutes) Register(r *gin.RouterGroup) {
	urls := r.Group("/urls")
	urls.Use(middleware.AuthMiddleware())
	{
		urls.GET("/:short_key", urlRoute.handler.DecryptShortKey)
		urls.POST("", urlRoute.handler.CreateShortUrl)
	}
}
