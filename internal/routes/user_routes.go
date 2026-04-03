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
		users.GET("/:uuid", userRoute.handler.GetUserByUUID)

		users.POST("/create", userRoute.handler.Create)
		users.PUT("/:uuid", userRoute.handler.Update)
		users.PUT("/:uuid/restore", userRoute.handler.RestoreUser)

		users.DELETE("/:uuid", userRoute.handler.SoftDelteUser)
		users.DELETE("/:uuid/clean", userRoute.handler.DeleteUser)
	}
}
