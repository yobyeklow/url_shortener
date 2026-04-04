package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var stackLineRegres = regexp.MustCompile(`(.+\.go:\d+)`)

func ExtractFirstAppStackLine(stack []byte) string {
	lines := bytes.Split(stack, []byte("\n"))
	var cleanLine string
	for _, line := range lines {
		if bytes.Contains(line, []byte(".go")) &&
			!bytes.Contains(line, []byte("/runtime/")) &&
			!bytes.Contains(line, []byte("/debug/")) &&
			!bytes.Contains(line, []byte("recovery_middleware.go")) {
			cleanLine = strings.TrimSpace(string(line))
			match := stackLineRegres.FindStringSubmatch(cleanLine)
			if len(match) > 1 {
				return match[1]
			}

		}
	}
	return ""
}
func RecoveryMiddleWare(recoveryLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				stack_at := ExtractFirstAppStackLine(stack)
				recoveryLogger.Error().
					Str("path", ctx.Request.URL.Path).
					Str("method", ctx.Request.Method).
					Str("client_ip", ctx.ClientIP()).
					Str("panic", fmt.Sprintf("%v", err)).
					Str("stack_at", stack_at).
					Str("stack", string(stack)).
					Msg("Panic occurred!")
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"Massage": "Try it again later!",
				})
			}
		}()
		ctx.Next()
	}
}
