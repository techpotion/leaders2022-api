package database

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/techpotion/leadershack2022/api/config"
	"go.uber.org/zap"
)

type Postgres struct {
	Pool         *pgxpool.Pool
	MetricsTimer time.Ticker
}

type logger struct {
	pool *pgxpool.Pool
	mx   sync.Mutex
}

func (l *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	z := zap.S().With("context", "postgres_logger", "msg", msg, "data", data)

	switch level {
	case pgx.LogLevelError:
		z.Error("Postgres error")
	case pgx.LogLevelWarn:
		z.Warn("Postgres warning")
	}
}

func (l *logger) setPool(pool *pgxpool.Pool) {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.pool = pool
}

func NewPostgresDatabase(cfg *config.Config) (*Postgres, error) {
	urlBuilder := new(strings.Builder)
	urlBuilder.WriteString(
		fmt.Sprintf("%s?pool_max_conns=%d&pool_max_conn_idle_time=30s&sslmode=disable&statement_cache_mode=describe",
			cfg.DBPgURI,
			cfg.MaxOpenConn,
		),
	)

	pgcfg, err := pgxpool.ParseConfig(urlBuilder.String())
	if err != nil {
		return nil, err
	}

	l := &logger{}
	pgcfg.ConnConfig.Logger = l
	pgcfg.ConnConfig.LogLevel = pgx.LogLevelInfo

	pool, err := pgxpool.ConnectConfig(context.Background(), pgcfg)
	if err != nil {
		return nil, err
	}

	l.setPool(pool)

	db := Postgres{
		Pool: pool,
	}

	return &db, nil
}
