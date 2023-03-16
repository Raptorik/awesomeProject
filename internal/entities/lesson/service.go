package lesson

import (
	"awesomeProject/pkg/logging"
	"context"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) FindAll(ctx context.Context) ([]Lesson, error) {
	return s.repository.FindAll(ctx)
}

func (s *Service) FindOne(ctx context.Context, id string) (Lesson, error) {
	return s.repository.FindOne(ctx, id)
}

func (s *Service) Update(ctx context.Context, lesson *Lesson) error {
	return s.repository.Update(ctx, lesson)
}
