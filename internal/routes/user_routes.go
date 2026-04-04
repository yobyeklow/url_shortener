package routes

import (
	"url_shortener/internal/handler"
	"url_shortener/internal/middleware"

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
	users.POST("/create", userRoute.handler.Create)
	private := users.Group("")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/:uuid", userRoute.handler.GetUserByUUID)
		private.GET("/", userRoute.handler.GetAllUser)
		private.GET("/soft-delete", userRoute.handler.GetSoftDeleteUsers)

		private.PUT("/:uuid", userRoute.handler.Update)
		private.PUT("/:uuid/restore", userRoute.handler.RestoreUser)

		private.DELETE("/:uuid", userRoute.handler.SoftDelteUser)
		private.DELETE("/:uuid/clean", userRoute.handler.DeleteUser)
	}
}
