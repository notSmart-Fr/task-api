package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	response "task-api/internal/platform/http"
)

type UpdateTaskRequest struct {
	ID    int     `json:"-"`
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

type UpdateTaskEndpoint struct {
	store Store
}

func NewUpdateTaskEndpoint(store Store) *UpdateTaskEndpoint {
	return &UpdateTaskEndpoint{store: store}
}

// @Summary Update task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param input body UpdateTaskRequest true "Update Payload"
// @Success 200 {object} Task
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /tasks/{id} [put]
func (e *UpdateTaskEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	req.ID = id

	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		response.Error(w, http.StatusBadRequest, "Task title cannot be empty")
		return
	}

	task, err := e.store.Update(req.ID, req.Title, req.Done)
	if err == ErrNotFound {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Task %d not found", req.ID))
		return
	} else if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	response.JSON(w, http.StatusOK, task)
}
