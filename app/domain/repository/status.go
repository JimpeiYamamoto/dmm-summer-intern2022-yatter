package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	PostStatus(ctx context.Context, entity *object.Status) error
	FindById(ctx context.Context, id int64) (*object.Status, error)
}
