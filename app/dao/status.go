package dao

import (
	"context"
	"database/sql"
	"errors"
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
		return err
	}
	entity.ID, _ = result.LastInsertId()
	return nil
}

func debugTable(ctx context.Context, s *status) {
	rows, _ := s.db.QueryContext(ctx, "SELECT  account_id, content FROM STATUS")
	defer rows.Close()

	var a object.Status

	for rows.Next() {
		rows.Scan(&a.AccountID, &a.Content)
		fmt.Println(a.AccountID, a.Content)
	}
}

func (r *status) GetTimelinesPublic(ctx context.Context, q object.Query) ([]object.Status, error) {
	msg := "SELECT id, account_id, content FROM status WHERE "
	msg += "id > " + q.SinceID
	msg += " AND id < " + q.MaxID
	l, err := strconv.Atoi(q.Limit)
	if err != nil {
		fmt.Println(fmt.Errorf("err: %w", err))
	}
	if l > 80 {
		q.Limit = "80"
	} else if q.Limit == "" {
		q.Limit = "40"
	}
	msg += " LIMIT " + q.Limit
	fmt.Println("================")
	fmt.Println(msg)
	fmt.Println("================")
	rows, err := r.db.QueryContext(ctx,
		msg,
	)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	a := make([]object.Status, 0)
	var s object.Status
	for rows.Next() {
		rows.Scan(&s.ID, &s.AccountID, &s.Content)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}
