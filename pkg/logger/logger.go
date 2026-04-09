package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type contextKey string

const TraceIdKey contextKey = "trace_id"

var Log *zerolog.Logger

type LoggerConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      string
	IsDev      string
}
type PrettyJSONWriter struct {
	Writer io.Writer
}

func (w PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJson bytes.Buffer

	if err := json.Indent(&prettyJson, p, "", "  "); err != nil {
		return w.Writer.Write(p)
	}
	return w.Writer.Write(prettyJson.Bytes())
}

func NewLogger(config LoggerConfig) *zerolog.Logger {
	zerolog.TimestampFieldName = time.RFC3339
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	var writter io.Writer
	if config.IsDev == "develope" {
		writter = PrettyJSONWriter{Writer: os.Stdout}
	} else {
		writter = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxAge,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
	}
	logger := zerolog.New(writter).With().Timestamp().Logger()
	return &logger
}
func GetTraceID(ctx context.Context) string {
	if traceId, ok := ctx.Value(TraceIdKey).(string); ok {
		return traceId
	}
	return ""
}
