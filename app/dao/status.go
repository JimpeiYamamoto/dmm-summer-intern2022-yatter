package dao

import (
	"context"
	"fmt"
	"strconv"
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

func (r *status) GetTimelinesPublic(ctx context.Context, q object.Query) ([]object.Status, error) {
	msg := "SELECT id, account_id, content, create_at FROM status WHERE "
	msg += "id > " + q.SinceID
	msg += " AND id < " + q.MaxID
	l, err := strconv.Atoi(q.Limit)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if l > 80 {
		q.Limit = "80"
	} else if q.Limit == "" {
		q.Limit = "40"
	}
	msg += " LIMIT " + q.Limit
	rows, err := r.db.QueryContext(ctx,
		msg,
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

func (r *status) FindById(ctx context.Context, id int64) (*object.Status, error) {
	entity := new(object.Status)
	err := r.db.QueryRowxContext(
		ctx,
		"select * from status where id = ?",
		id,
	).StructScan(entity)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}
