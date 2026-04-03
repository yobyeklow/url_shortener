package app

import (
	"url_shortener/internal/handler"
	"url_shortener/internal/repository"
	"url_shortener/internal/routes"
	"url_shortener/internal/services"
	"url_shortener/pkg/auth"
	"url_shortener/pkg/cache"
)

type AuthModule struct {
	route routes.Routes
}

func NewAuthModule(ctx *ModuleContext, tokenService auth.TokenService, cache cache.RedisService) *AuthModule {
	authRepo := repository.NewSQLUserRepository(ctx.db)
	authService := services.NewAuthServices(authRepo, tokenService, cache)
	authHandler := handler.NewAuthHandler(authService)
	authRoutes := routes.NewAuthRoute(authHandler)
	return &AuthModule{
		route: authRoutes,
	}
}
func (module *AuthModule) Routes() routes.Routes {
	return module.route
}
