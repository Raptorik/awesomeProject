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
func (s *Service) TranslateLessonName(ctx context.Context, lessonID string) error {
	lesson, err := s.repository.FindOne(ctx, lessonID)
	if err != nil {
		return err
	}

	if lesson.Language == "en" {
		// no need to translate
		return nil
	}

	// translate lesson name
	translatedName, err := s.repository.TranslateLessonName(ctx, lesson.Name, lesson.Language)
	if err != nil {
		return err
	}

	// update repository with translated name
	lesson.TranslatedName = translatedName
	if err := s.repository.Update(ctx, &lesson); err != nil {
		return err
	}

	return nil
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
