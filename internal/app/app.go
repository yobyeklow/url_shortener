package app

import (
	"context"
	"log"
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
		log.Fatalf("Validator init failed %v", err)
	}
	r := gin.Default()
	//Connect DB
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database init failed %v", err)
	}
	redisClient := config.NewRedisClient()
	cacheRedis := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedis)
	ctx := &ModuleContext{
		db: database.DB,
	}
	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, cacheRedis),
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
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server %v", err)
		}
	}()
	<-quitSrv
	log.Println("Shutdown signal received!")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server shutdown!")
	return nil
}
func getModulesRoute(modules []Module) []routes.Routes {
	routeList := make([]routes.Routes, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}
