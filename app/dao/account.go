package dao

import (
	"context"
	"database/sql"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// Create new user account
func (r *account) CreateNewAccount(ctx context.Context, entity object.Account) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO account (id, username, password_hash) VALUES (?, ?, ?)",
		entity.ID,
		entity.Username,
		entity.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// Fetch user by ID
func (r *account) FindByUserID(ctx context.Context, id int64) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(
		ctx,
		"select * from account where id = ?",
		id,
	).StructScan(entity)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	entity.FollowingCount, entity.FollowerCount, err = r.CountFollwingFollower(ctx, entity.ID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

// Fetch account which has specified username
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	entity.FollowingCount, entity.FollowerCount, err = r.CountFollwingFollower(ctx, entity.ID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

// Follow user
func (r *account) FollowUser(ctx context.Context, myId, targetId int64) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO relation (follower_id, followee_id) VALUES (?, ?)",
		myId,
		targetId,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func QueryFollowingUser(ctx context.Context, r *account, id int64, limit string) (*sql.Rows, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"select followee_id from relation where follower_id = ? LIMIT ?",
		id,
		limit,
	)
	if err != nil {
		return nil, nil
	}
	return rows, err
}

// Fetch following user account
func (r *account) GetFollowings(ctx context.Context, username string, limit string) ([]object.Account, error) {
	a, err := r.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	rows, err := QueryFollowingUser(ctx, r, a.ID, limit)
	if err != nil {
		return nil, nil
	}
	as := make([]object.Account, 0)
	for rows.Next() {
		var followeeId int64
		rows.Scan(&followeeId)
		a, err := r.FindByUserID(ctx, followeeId)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		as = append(as, *a)
	}
	return as, nil
}

func QueryFollowerUser(ctx context.Context, r *account, id int64, q object.Query) (*sql.Rows, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT follower_id FROM relation WHERE followee_id = ? AND follower_id < ? AND follower_id > ? LIMIT ?",
		id,
		q.MaxID,
		q.SinceID,
		q.Limit,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Fetch follower user account
func (r *account) GetFollowers(ctx context.Context, username string, query object.Query) ([]object.Account, error) {
	a, err := r.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	rows, err := QueryFollowerUser(ctx, r, a.ID, query)
	if err != nil {
		return nil, nil
	}
	as := make([]object.Account, 0)
	for rows.Next() {
		var followeeId int64
		rows.Scan(&followeeId)
		a, err := r.FindByUserID(ctx, followeeId)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		as = append(as, *a)
	}
	return as, nil
}

// Unfollow user
func (r *account) UnfollowUser(ctx context.Context, myId, targetId int64) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM relation WHERE (follower_id, followee_id) = (?, ?)",
		targetId,
		myId,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// Fetch relation ships auth account with request account
func (r *account) GetRelationships(ctx context.Context, myId, targetId int64) (*object.Relation, error) {
	var entity object.Relation
	entity.Id = targetId
	rows, _ := r.db.QueryContext(
		ctx,
		"select * from relation where followee_id = ? AND follower_id = ?",
		myId,
		targetId,
	)
	if rows.Next() {
		entity.FollowedBy = true
	}
	rows, _ = r.db.QueryContext(
		ctx,
		"select * from relation where followee_id = ? AND follower_id = ?",
		targetId,
		myId,
	)
	if rows.Next() {
		entity.Following = true
	}
	return &entity, nil
}

// Count follower and following account
func (r *account) CountFollwingFollower(ctx context.Context, id int64) (int, int, error) {
	q := object.Query{
		OnlyMedia: "",
		MaxID:     "100000",
		SinceID:   "0",
		Limit:     "80",
	}
	rows, err := QueryFollowingUser(ctx, r, id, q.Limit)
	if err != nil {
		return 0, 0, err
	}
	n1 := 0
	for rows.Next() {
		n1++
	}
	rows, err = QueryFollowerUser(ctx, r, id, q)
	if err != nil {
		return 0, 0, err
	}
	n2 := 0
	for rows.Next() {
		n2++
	}
	return n1, n2, nil
}
