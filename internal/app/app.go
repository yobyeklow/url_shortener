package app

import (
	"log"
	"url_shortener/internal/config"
	"url_shortener/internal/routes"
	"url_shortener/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

type Module interface {
	Routes() routes.Routes
}
type Application struct {
	cfg     *config.Config
	router  *gin.Engine
	modules []Module
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validator init failed %v", err)
	}
	r := gin.Default()
	loadEnv()
	modules := []Module{
		NewUserModule(),
	}
	routes.RegisterRoutes(r, getModuleRoute(modules)...)
	return &Application{
		router:  r,
		cfg:     cfg,
		modules: modules,
	}
}
func (app *Application) Run() error {
	return app.router.Run(app.cfg.ServerAddress)
}
func getModuleRoute(modules []Module) []routes.Routes {
	routeList := make([]routes.Routes, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}
func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
