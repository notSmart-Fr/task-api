package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	response "task-api/internal/platform/http"
)

type GetTaskRequest struct {
	ID int
}

type GetTaskEndpoint struct {
	store Store
}

func NewGetTaskEndpoint(store Store) *GetTaskEndpoint {
	return &GetTaskEndpoint{store: store}
}

// @Summary Get task by ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} Task
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /tasks/{id} [get]
func (e *GetTaskEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	req := GetTaskRequest{ID: id}
	task, err := e.store.GetByID(req.ID)
	if err == ErrNotFound {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Task %d not found", req.ID))
		return
	} else if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	response.JSON(w, http.StatusOK, task)
}
