package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url_shortener/internal/config"
	"url_shortener/internal/database"
	"url_shortener/internal/database/sqlc"
	"url_shortener/internal/routes"
	"url_shortener/internal/validation"
	"url_shortener/pkg/auth"
	"url_shortener/pkg/cache"
	"url_shortener/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() routes.Routes
}
type Application struct {
	cfg     *config.Config
	router  *gin.Engine
	modules []Module
}
type ModuleContext struct {
	db sqlc.Querier
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Validator init failed %v")
	}
	r := gin.Default()
	//Connect DB
	if err := database.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database init failed %v")
	}
	redisClient := config.NewRedisClient()
	cacheRedis := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedis)
	ctx := &ModuleContext{
		db: database.DB,
	}
	modules := []Module{
		NewUserModule(ctx, redisClient),
		NewAuthModule(ctx, tokenService, cacheRedis),
		NewUrlModule(ctx, redisClient),
	}

	routes.RegisterRoutes(r, tokenService, cacheRedis, getModulesRoute(modules)...)
	return &Application{
		router:  r,
		cfg:     cfg,
		modules: modules,
	}
}
func (app *Application) Run() error {
	server := &http.Server{
		Addr:    app.cfg.Address,
		Handler: app.router,
	}
	quitSrv := make(chan os.Signal, 1)
	// syscall.SIGNINT -> Ctrl + C -> End  Process
	// syscall.SIGTERM -> Kill Service
	// syscall.SIGHUP -> Reload Service
	signal.Notify(quitSrv, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	logger.Log.Info().Msgf("Server is running at %s", app.cfg.Address)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Failed to start server %v")
		}
	}()
	<-quitSrv
	logger.Log.Info().Msg("Shutdown signal received!")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server forced to shutdown: %v")
	}
	logger.Log.Info().Msg("Server shutdown!")
	return nil
}
func getModulesRoute(modules []Module) []routes.Routes {
	routeList := make([]routes.Routes, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}
