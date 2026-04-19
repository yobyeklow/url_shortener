package database

import (
	"context"
	"fmt"
	"time"
	"url_shortener/internal/config"
	"url_shortener/internal/database/sqlc"
	"url_shortener/internal/utils"
	"url_shortener/pkg/logger"
	"url_shortener/pkg/pgx"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

var DBPool *pgxpool.Pool
var DB sqlc.Querier

func InitDB() error {
	connectStr := config.NewConfig().DNS()

	sqlLogger := utils.NewLoggerWithPath("sql.log", "info")
	conf, err := pgxpool.ParseConfig(connectStr)

	if err != nil {
		return fmt.Errorf("Error parsing DB config: %v", err)
	}

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: &pgx.PgxZerologTracer{
			Logger:         *sqlLogger,
			SlowQueryLimit: 500 * time.Microsecond,
		},
		LogLevel: tracelog.LogLevelDebug,
	}

	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	DBPool, err = pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("Error creating Database poo: %v", err)
	}
	DB = sqlc.New(DBPool)
	if err := DBPool.Ping(ctx); err != nil {
		return fmt.Errorf("DB Ping error: %v", err)
	}

	logger.Log.Info().Msg("Connected Database successfully!")
	return nil
}
