package lesson

import "context"

type Repository interface {
	FindAll(ctx context.Context) (u []Lesson, err error)
}
