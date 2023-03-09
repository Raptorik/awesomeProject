package lesson

import (
	"awesomeProject/internal/apperror"
	"awesomeProject/internal/handlers"
	"awesomeProject/pkg/logging"
	translation_lesson "awesomeProject/pkg/translation"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, translator *translation_lesson.GoogleTranslator, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/lessons", apperror.Middleware(h.GetAllLessons))
	router.HandlerFunc(http.MethodGet, "/lessons/:id", apperror.Middleware(h.GetLesson))
	router.HandlerFunc(http.MethodPatch, "/lessons/:lang", apperror.Middleware(h.UpdateLesson))
}

func (h *handler) GetAllLessons(w http.ResponseWriter, r *http.Request) error {
	lessons, err := h.repository.FindAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get lessons: %v", err)))
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lessons); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to encode response: %v", err)))
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *handler) GetLesson(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	lessonID := params.ByName("uuid")

	lesson, err := h.repository.FindOne(r.Context(), lessonID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get lesson: %v", err)))
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lesson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to encode response: %v", err)))
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
func (h *handler) UpdateLesson(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	lessonID := chi.URLParam(r, "id")

	lesson, err := h.repository.FindOne(ctx, lessonID)
	if err != nil {
		h.logger.Errorf("failed to fetch lesson with ID %s: %v", lessonID, err)
		return fmt.Errorf("lesson not found")
	}

	googleTranslator, err := translation_lesson.NewGoogleTranslator("YOUR_API_KEY_HERE")
	if err != nil {
		h.logger.Errorf("failed to create Google translator: %v", err)
		return fmt.Errorf("failed to create Google translator")
	}

	err = translation_lesson.TranslateLessonName(lesson, "ukr", googleTranslator)
	if err != nil {
		h.logger.Errorf("failed to translate lesson name: %v", err)
		return fmt.Errorf("failed to translate lesson name")
	}

	var updatedLesson Lesson
	err = json.NewDecoder(r.Body).Decode(&updatedLesson)
	if err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		return fmt.Errorf("invalid request payload")
	}

	lesson.Name = updatedLesson.Name
	lesson.Courses = updatedLesson.Courses

	err = h.repository.Update(ctx, &lesson)
	if err != nil {
		h.logger.Errorf("failed to update lesson with ID %s: %v", lessonID, err)
		return fmt.Errorf("failed to update lesson")
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
