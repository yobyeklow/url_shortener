package routes

import (
	"url_shortener/internal/middleware"
	"url_shortener/internal/utils"
	"url_shortener/pkg/auth"
	"url_shortener/pkg/cache"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, authService auth.TokenService, cacheService cache.RedisService, routes ...Routes) {
	httpLogger := utils.NewLoggerWithPath("http.log", "info")
	recoveryLogger := utils.NewLoggerWithPath("recovery.log", "warning")
	rateLimitLogger := utils.NewLoggerWithPath("ratelimit.log", "warning")

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(
		middleware.RateLimiterMiddleware(rateLimitLogger),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleWare(recoveryLogger),
	)

	api := r.Group("/api/v1")
	middleware.InitAuthMiddlware(authService, cacheService)

	for _, route := range routes {
		route.Register(api)
	}
	//Throw the error if use wrong method for the routes
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"Error": "Not Found",
			"path":  ctx.Request.URL.Path,
		})
	})
}
