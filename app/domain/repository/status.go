package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	PostStatus(ctx context.Context, entity object.Status) error
}
