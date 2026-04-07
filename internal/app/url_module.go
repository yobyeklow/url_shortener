package app

import (
	"url_shortener/internal/handler"
	"url_shortener/internal/repository"
	"url_shortener/internal/routes"
	"url_shortener/internal/services"

	"github.com/redis/go-redis/v9"
)

type UrlModule struct {
	route routes.Routes
}

func NewUrlModule(ctx *ModuleContext, redisClient *redis.Client) *UrlModule {
	urlRepo := repository.NewSQLUrlRepository(ctx.db)
	urlService := services.NewUrlService(urlRepo, redisClient)
	urlHandler := handler.NewUrlHandler(urlService)
	urlRoutes := routes.NewUrlRoutes(urlHandler)
	return &UrlModule{
		route: urlRoutes,
	}
}
func (module *UrlModule) Routes() routes.Routes {
	return module.route
}
