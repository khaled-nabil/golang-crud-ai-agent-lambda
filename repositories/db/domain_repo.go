package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	DomainRepository struct {
		agent *pgxpool.Pool
	}
)

const (
	domainTable = "domain"
)

func NewDomainRepository(db *PostgresRepo) *DomainRepository {
	return &DomainRepository{db.agent}
}

func (b *DomainRepository) GetDomainInstructions(name string) (string, error) {
	var ins string

	query := fmt.Sprintf("SELECT instructions FROM %s WHERE name = $1", domainTable)

	err := b.agent.QueryRow(context.Background(), query, name).Scan(&ins)
	if err != nil {
		return "", fmt.Errorf("query domain instructions: %w", err)
	}

	return ins, nil
}
