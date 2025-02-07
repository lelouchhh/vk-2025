package repository

import (
	"backend/domain"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type PostgresRepository interface {
	SavePingResult(ctx context.Context, result domain.PingResult) error
	GetAllContainers(ctx context.Context) ([]domain.Container, error)
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) PostgresRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) SavePingResult(ctx context.Context, result domain.PingResult) error {
	query := `
        INSERT INTO containers (ip_address, last_ping,ping_time, status)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.db.ExecContext(ctx, query, result.IP, result.LastSuccess, result.PingTime, result.Status)
	if err != nil {
		log.Printf("Failed to save ping result: %v", err)
		return err
	}
	return nil
}

func (r *postgresRepository) GetAllContainers(ctx context.Context) ([]domain.Container, error) {
	query := `
        SELECT id, ip_address, last_ping, status, ping_time
        FROM containers
    `
	var containers []domain.Container
	err := r.db.SelectContext(ctx, &containers, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Failed to fetch containers: %v", err)
		return nil, err
	}
	return containers, nil
}
