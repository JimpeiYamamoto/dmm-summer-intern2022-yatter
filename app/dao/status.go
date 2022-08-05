package dao

import (
	"context"
	"fmt"
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

// Post authentication user new status
func (r *status) PostStatus(ctx context.Context, entity *object.Status) error {
	result, err := r.db.ExecContext(
		ctx,
		"INSERT INTO status (account_id, content) VALUES (?, ?)",
		entity.AccountID,
		entity.Content,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	entity.ID, _ = result.LastInsertId()
	o, err := r.FindById(ctx, entity.ID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	entity.CreateAt = o.CreateAt
	return nil
}

// Fetch Status by ID
func (r *status) FindById(ctx context.Context, id int64) (*object.Status, error) {
	entity := new(object.Status)
	err := r.db.QueryRowxContext(
		ctx,
		"SELECT * FROM status WHERE id = ?",
		id,
	).StructScan(entity)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

// Delete authentication user status
func (r *status) DeleteStatus(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM status WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// Fetch timeline all user
func (r *status) GetTimelineHome(ctx context.Context, q object.Query, accountID int64) ([]object.Status, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, account_id, content, create_at 
		FROM status 
		WHERE account_id = ?  
		AND id > ? 
		AND id < ? 
		LIMIT ?`,
		accountID,
		q.SinceID,
		q.MaxID,
		q.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()
	ss := make([]object.Status, 0)
	var s object.Status
	for rows.Next() {
		rows.Scan(&s.ID, &s.AccountID, &s.Content, &s.CreateAt)
		ss = append(ss, s)
	}
	return ss, nil
}

// Fetch timeline following user
func (r *status) GetTimelinesPublic(ctx context.Context, q object.Query) ([]object.Status, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, account_id, content, create_at
		FROM status
		WHERE id > ?
		AND id < ?
		LIMIT ?`,
		q.SinceID,
		q.MaxID,
		q.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()
	a := make([]object.Status, 0)
	var s object.Status
	for rows.Next() {
		rows.Scan(&s.ID, &s.AccountID, &s.Content, &s.CreateAt)
		a = append(a, s)
	}
	return a, nil
}
