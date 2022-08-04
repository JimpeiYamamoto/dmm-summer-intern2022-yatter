package dao

import (
	"context"
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

func DebugTable(ctx context.Context, r *account) {
	rows, _ := r.db.QueryContext(
		ctx,
		"SELECT id, username FROM account",
	)
	var a object.Account
	for rows.Next() {
		rows.Scan(&a.ID, &a.Username)
		fmt.Println(a)
	}
}

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

func debugRelation(ctx context.Context, r *account) {
	rows, _ := r.db.QueryContext(ctx, "SELECT * FROM relation")
	defer rows.Close()

	for rows.Next() {
		var follower_id, followee_id int64
		rows.Scan(&follower_id, &followee_id)
		res, _ := r.FindByUserID(ctx, follower_id)
		follower_name := res.Username

		res, _ = r.FindByUserID(ctx, followee_id)
		followee_name := res.Username

		fmt.Printf("%s => %s\n", follower_name, followee_name)
	}
}

func (r *account) GetFollowingUser(ctx context.Context, username string, limit string) ([]object.Account, error) {
	a, err := r.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	rows, err := r.db.QueryContext(
		ctx,
		"select followee_id from relation where follower_id = ? LIMIT ?",
		a.ID,
		limit,
	)
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

func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

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
	return entity, nil
}
