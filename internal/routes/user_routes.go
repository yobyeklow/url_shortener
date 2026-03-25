package routes

import (
	"url_shortener/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handler.UserHandler
}

func NewUserRoutes(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}

func (userRoute *UserRoutes) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("/", userRoute.handler.GetAllUser)
		users.POST("/create", userRoute.handler.GetAllUser)
		users.DELETE("/delete", userRoute.handler.GetAllUser)
		users.PATCH("/update", userRoute.handler.GetAllUser)
	}
}
