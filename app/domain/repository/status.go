package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Post authentication user new status
	PostStatus(ctx context.Context, entity *object.Status) error
	// Fetch Status by ID
	FindById(ctx context.Context, id int64) (*object.Status, error)
	// Delete authentication user status
	DeleteStatus(ctx context.Context, id string) error
	// Fetch timeline all user
	GetTimelinesPublic(ctx context.Context, q object.Query) ([]object.Status, error)
	// Fetch timeline following user
	GetTimelineHome(ctx context.Context, q object.Query, accountID int64) ([]object.Status, error)
}
