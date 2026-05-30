package provider

import (
	"context"
	"eventhandler/util"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := util.Configuration.Postgres

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		strings.Join(cfg.Options, "&"),
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	if cfg.MaxConns > 0 {
		poolConfig.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolConfig.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetimeSecs > 0 {
		poolConfig.MaxConnLifetime = time.Duration(cfg.MaxConnLifetimeSecs) * time.Second
	}
	if cfg.MaxConnIdleTimeSecs > 0 {
		poolConfig.MaxConnIdleTime = time.Duration(cfg.MaxConnIdleTimeSecs) * time.Second
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
