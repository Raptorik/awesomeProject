package course

import (
	"awesomeProject/pkg/logging"
	"context"
	"fmt"
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
func (s *Service) CreateCourse(ctx context.Context, c *Course) error {
	if c.Name == "" {
		return fmt.Errorf("course name cannot be empty")
	}
	if err := s.repository.Create(ctx, c); err != nil {
		return fmt.Errorf("faile to create course due to errror: %v", err)
	}
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]Course, error) {
	all, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all course due to error: %v", err)
	}
	return all, nil
}
func (s *Service) Update(ctx context.Context, course *Course) error {
	_, err := s.repository.FindOne(ctx, course.ID)
	if err != nil {
		return fmt.Errorf("failed to update course: %v", err)
	}
	err = s.repository.Update(ctx, course)
	if err != nil {
		return fmt.Errorf("failed to update course: %v", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("course ID cannot be empty")
	}
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete course with ID %s due to error: %v", id, err)
	}
	return nil
}
