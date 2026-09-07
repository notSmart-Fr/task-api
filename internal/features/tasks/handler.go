package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	response "task-api/internal/platform/http"
)

type Handler struct {
	store Store
}

// RegisterHandlers registers task endpoints on the provided ServeMux.
func RegisterHandlers(mux *http.ServeMux, store Store) {
	h := &Handler{store: store}

	mux.HandleFunc("GET /tasks", h.list)
	mux.HandleFunc("GET /tasks/{id}", h.get)
	mux.HandleFunc("POST /tasks", h.create)
	mux.HandleFunc("PUT /tasks/{id}", h.update)    // Added Stage 4 route
	mux.HandleFunc("DELETE /tasks/{id}", h.delete) // Added Stage 4 route
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.store.GetAll())
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Task %d not found", id))
		return
	}

	response.JSON(w, http.StatusOK, task)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var input CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		response.Error(w, http.StatusBadRequest, "Task title is required and cannot be empty")
		return
	}

	task := h.store.Create(input.Title)
	response.JSON(w, http.StatusCreated, task)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	var input UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		response.Error(w, http.StatusBadRequest, "Task title cannot be empty")
		return
	}

	task, err := h.store.Update(id, input)
	if err != nil {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Task %d not found", id))
		return
	}

	response.JSON(w, http.StatusOK, task)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	if err := h.store.Delete(id); err != nil {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Task %d not found", id))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
