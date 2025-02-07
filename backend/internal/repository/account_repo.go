package repository

import (
	"backend/domain"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account domain.Account) error
	GetAccountByLogin(ctx context.Context, login string) (*domain.Account, error)
}

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account domain.Account) error {
	query := `
        INSERT INTO account (login, password)
        VALUES ($1, $2)
    `
	_, err := r.db.ExecContext(ctx, query, account.Login, account.Password)
	if err != nil {
		log.Printf("Failed to create account: %v", err)
		return err
	}
	return nil
}

func (r *accountRepository) GetAccountByLogin(ctx context.Context, login string) (*domain.Account, error) {
	query := `
        SELECT id, login, password
        FROM account
        WHERE login = $1
    `
	var account domain.Account
	err := r.db.GetContext(ctx, &account, query, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Failed to fetch account by login: %v", err)
		return nil, err
	}
	return &account, nil
}
