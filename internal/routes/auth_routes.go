package routes

import (
	"url_shortener/internal/handler"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	handler *handler.AuthHandler
}

func NewAuthRoute(handler *handler.AuthHandler) *AuthRoute {
	return &AuthRoute{
		handler: handler,
	}
}
func (authRoute *AuthRoute) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", authRoute.handler.Login)
		auth.POST("/logout", authRoute.handler.Logout)
		auth.POST("/refresh-token", authRoute.handler.RefreshToken)
	}
}
