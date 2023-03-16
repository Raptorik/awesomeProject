package lesson

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context) (lesson []Lesson, err error)
	FindOne(ctx context.Context, id string) (Lesson, error)
	Update(ctx context.Context, lesson *Lesson) error
}
