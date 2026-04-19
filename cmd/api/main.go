package main

import (
	"log"
	"os"
	"path/filepath"
	"url_shortener/internal/app"
	"url_shortener/internal/config"
	"url_shortener/internal/utils"
	"url_shortener/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	rootDir := mustGetWorkingDir()
	logFile := filepath.Join(rootDir, "internal/logs/app.log")
	logger.InitLogger(logger.LoggerConfig{
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		Level:      "info",
		IsDev:      utils.GetEnv("APP_EVN", "develope"),
	})
	loadEnv(filepath.Join(rootDir, "env"))
	cfg := config.NewConfig()
	app := app.NewApplication(cfg)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func mustGetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get working dir:", err)
	}
	return dir
}
func loadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		logger.Log.Warn().Msg("No .env file found!")
	} else {
		logger.Log.Info().Msg("Loaded .ENV file successfully!")
	}
}
