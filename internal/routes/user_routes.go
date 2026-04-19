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
	//Public Route
	users.POST("/create", userRoute.handler.Create)
	auth := users.Group("")
	auth.Use(middleware.AuthMiddleware())
	{

		auth.GET("/:uuid", userRoute.handler.GetUserByUUID)
		auth.GET("/", userRoute.handler.GetAllUser)
		auth.PUT("/:uuid", userRoute.handler.Update)

		adminMod := auth.Group("")
		//User Role: 1 - User, 2 - Moderator, 3 - Adminstator
		adminMod.Use(middleware.RequirePermission("Administrator", "Moderator"))
		{
			adminMod.GET("/soft-delete", userRoute.handler.GetSoftDeleteUsers)
			adminMod.PUT("/:uuid/restore", userRoute.handler.RestoreUser)
			adminMod.DELETE("/:uuid", userRoute.handler.SoftDelteUser)
			adminMod.DELETE("/:uuid/clean", userRoute.handler.DeleteUser)
		}
	}
}
