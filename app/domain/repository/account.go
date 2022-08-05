package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Account interface {
	// Create new user account
	CreateNewAccount(ctx context.Context, entity object.Account) error
	// Fetch user by ID
	FindByUserID(ctx context.Context, id int64) (*object.Account, error)
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)
	// Follow user
	FollowUser(ctx context.Context, myId, targetId int64) error
	// Fetch following user account
	GetFollowings(ctx context.Context, username, limit string) ([]object.Account, error)
	// Fetch follower user account
	GetFollowers(ctx context.Context, username string, query object.Query) ([]object.Account, error)
	// Unfollow user
	UnfollowUser(ctx context.Context, myId, targetId int64) error
	// Fetch relation ships auth account with request account
	GetRelationships(ctx context.Context, myId, targetId int64) (*object.Relation, error)
	// Count follower and following account
	CountFollwingFollower(ctx context.Context, id int64) (int, int, error)
}
