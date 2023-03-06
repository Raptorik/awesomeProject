package lesson

import (
	"awesomeProject/internal/apperror"
	"awesomeProject/internal/handlers"
	"encoding/json"
	"fmt"

	"awesomeProject/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/lessons", apperror.Middleware(h.GetAllLessons))
	router.HandlerFunc(http.MethodGet, "/lessons/:id", apperror.Middleware(h.GetLesson))
	router.HandlerFunc(http.MethodPost, "/lessons/translate", apperror.Middleware(h.TranslateLessonName))
	router.HandlerFunc(http.MethodPatch, "/lessons/:id", apperror.Middleware(h.UpdateLesson))
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
func (h *handler) TranslateLessonName(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	lessonID := params.ByName("id")

	lesson, err := h.repository.FindOne(r.Context(), lessonID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to translate lesson name: %v", err)))
		return nil
	}
	if lesson.Language == "en" {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	translatedName, err := h.repository.TranslateLessonName(r.Context(), lesson.Name, lesson.Language)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to translate lesson name: %v", err)))
		return nil
	}
	lesson.TranslatedName = translatedName
	if err := h.repository.Update(r.Context(), &lesson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to update lesson: %v", err)))
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
	params := httprouter.ParamsFromContext(r.Context())
	lessonID := params.ByName("uuid")

	var lesson Lesson
	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to decode request body: %v", err)))
		return nil
	}

	lesson.ID = lessonID
	if err := h.repository.Update(r.Context(), &lesson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to update lesson: %v", err)))
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
