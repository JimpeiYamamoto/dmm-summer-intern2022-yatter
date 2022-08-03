package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	status struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

func (r *status) PostStatus(ctx context.Context, entity object.Status) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO status (id, account_id, content) VALUES (?, ?, ?)",
		entity.ID,
		entity.AccountID,
		entity.Content,
	)
	if err != nil {
		return err
	}
	return nil
}
