package utils

import (
	"fmt"
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
func NewLoggerWithPath(fileName string, level string) *zerolog.Logger {
	cwd, err := os.Getwd()
	fmt.Println(cwd)
	if err != nil {
		log.Fatal("❌ Unable to get working dir:", err)
	}
	logDir := filepath.Join(cwd, "..", "..", "internal/logs/", fileName)

	config := logger.LoggerConfig{
		Level:      level,
		Filename:   logDir,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      GetEnv("APP_EVN", "development"),
	}
	return logger.NewLogger(config)
}
