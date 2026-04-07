package app

import (
	"url_shortener/internal/handler"
	"url_shortener/internal/repository"
	"url_shortener/internal/routes"
	"url_shortener/internal/services"

	"github.com/redis/go-redis/v9"
)

type UserModule struct {
	route routes.Routes
}

func NewUserModule(ctx *ModuleContext, redisClient *redis.Client) *UserModule {
	userRepo := repository.NewSQLUserRepository(ctx.db)
	userService := services.NewUserService(userRepo, redisClient)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)
	return &UserModule{
		route: userRoutes,
	}
}
func (module *UserModule) Routes() routes.Routes {
	return module.route
}
