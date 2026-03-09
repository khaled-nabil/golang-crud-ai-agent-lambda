package db

import (
	"ai-agent/model/datamodels"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Repository struct {
		agent *pgxpool.Pool
		ctx   context.Context
	}
)

func New(config *datamodels.AppConfig) (*Repository, error) {
	ctx := context.Background()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{
		agent: pool,
		ctx:   ctx,
	}, nil
}

func (r *Repository) Close() error {
	if r.agent != nil {
		r.agent.Close()
	}
	return nil
}
