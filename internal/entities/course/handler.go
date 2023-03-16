package course

import (
	"awesomeProject/internal/apperror"
	"awesomeProject/internal/handlers"
	"awesomeProject/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	coursesURL = "/courses"
	courseURL  = "/courses/:uuid"
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
	router.HandlerFunc(http.MethodGet, coursesURL, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodPost, courseURL, apperror.Middleware(h.CreateCourse))
	router.HandlerFunc(http.MethodPut, courseURL, apperror.Middleware(h.Update))
	router.HandlerFunc(http.MethodDelete, courseURL, apperror.Middleware(h.Delete))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to get all courses: %v", err)))
		return nil
	}
	allBytes, err := json.Marshal(all)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to marshal courses: %v", err)))
		return nil
	}
	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)
	return nil
}

func (h *handler) CreateCourse(w http.ResponseWriter, r *http.Request) error {
	var c Course
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to decode course: %v", err)))
		return nil
	}
	if err := h.repository.Create(r.Context(), &c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to create course: %v", err)))
		return nil
	}
	cBytes, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to marshal course: %v", err)))
		return nil
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(cBytes)
	return nil
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	var c Course
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to decode course: %v", err)))
		return nil
	}
	c.ID = id
	if err := h.repository.Update(r.Context(), &c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to update course: %v", err)))
		return nil
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	if err := h.repository.Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to delete course: %v", err)))
		return nil
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
