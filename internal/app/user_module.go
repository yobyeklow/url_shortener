package app

import (
	"url_shortener/internal/handler"
	"url_shortener/internal/repository"
	"url_shortener/internal/routes"
	"url_shortener/internal/services"
)

type UserModule struct {
	route routes.Routes
}

func NewUserModule() *UserModule {
	userRepo := repository.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)

	return &UserModule{
		route: userRoutes,
	}
}
func (module *UserModule) Routes() routes.Routes {
	return module.route
}
