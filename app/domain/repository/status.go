package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	FindById(ctx context.Context, id int64) (*object.Status, error)
	PostStatus(ctx context.Context, entity *object.Status) error
	DeleteStatus(ctx context.Context, id string) error
	GetTimelinesPublic(ctx context.Context, q object.Query) ([]object.Status, error)
	GetTimelineHome(ctx context.Context, q object.Query, accountID int64) ([]object.Status, error)
}
