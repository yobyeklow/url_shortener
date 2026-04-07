package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"url_shortener/pkg/logger"

	"strconv"

	"github.com/rs/zerolog"
	"github.com/zeebo/xxh3"
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
func GenerateRandomKey(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b)
}
func Base62Encode(num int32) string {
	if num == 0 {
		return "0"
	}
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	for num > 0 {
		remainder := num % 62
		sb.WriteByte(base62Chars[remainder])
		num /= 62
	}
	//Reverse string
	runes := []rune(sb.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func GenerateHashedValue(
	iosDeepLink, iosFallbackURL, androidDeepLink, androidFallbackURL *string,
	defaultFallbackURL string,
) string {
	getStr := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	input := fmt.Sprintf("%s:%s:%s:%s:%s",
		getStr(iosDeepLink),
		getStr(iosFallbackURL),
		getStr(androidDeepLink),
		getStr(androidFallbackURL),
		defaultFallbackURL,
	)

	// Calculate 128-bit
	hash := xxh3.Hash128([]byte(input))

	// Convert 128-bit to 16 byte (little-endian)
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf[0:8], hash.Lo)
	binary.LittleEndian.PutUint64(buf[8:16], hash.Hi)

	return hex.EncodeToString(buf)
}
