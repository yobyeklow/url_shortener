package middleware

import (
	"context"
	"url_shortener/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		contextGo := context.WithValue(ctx.Request.Context(), logger.TraceIdKey, traceID)
		ctx.Request = ctx.Request.WithContext(contextGo) // Save traceID to context of Golang
		ctx.Writer.Header().Set("X-Trace-Id", traceID)   // Set TraceID to Response Header (Not Request Header)
		ctx.Set(string(logger.TraceIdKey), traceID)      // Save traceID to context of Gin
		ctx.Next()
	}
}
