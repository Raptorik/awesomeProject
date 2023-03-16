package course

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, course *Course) error
	FindAll(ctx context.Context) (_ []Course, err error)
	FindOne(ctx context.Context, id string) (Course, error)
	Update(ctx context.Context, course *Course) error
	Delete(ctx context.Context, id string) error
}
