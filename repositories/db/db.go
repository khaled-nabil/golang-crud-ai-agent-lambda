package db

import (
	"ai-agent/adapters/secrets"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	PostgresRepo struct {
		agent *pgxpool.Pool
		ctx   context.Context
	}
)

func NewPostgresRepo(config *secrets.AppConfig) (*PostgresRepo, error) {
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

	return &PostgresRepo{
		agent: pool,
		ctx:   ctx,
	}, nil
}

func (r *PostgresRepo) Close() error {
	if r.agent != nil {
		r.agent.Close()
	}
	return nil
}
