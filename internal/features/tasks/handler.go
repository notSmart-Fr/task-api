package tasks

import (
	"fmt"
	"net/http"
	"strconv"

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
