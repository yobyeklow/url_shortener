package utils

import (
	"log"
	"os"
	"path/filepath"
	"url_shortener/pkg/logger"

	"strconv"

	"github.com/rs/zerolog"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}
func NewLoggerWithPath(path string, level string) *zerolog.Logger {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get working dir:", err)
	}
	path = filepath.Join(cwd, "internal/logs/", path)
	config := logger.LoggerConfig{
		Filename:   path,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5, //days
		Compress:   true,
		Level:      level,
		IsDev:      GetEnv("APP_STATUS", "development"),
	}
	return logger.NewLogger(config)
}
