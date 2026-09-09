package tasks

import (
	"encoding/json"
	"net/http"
	"strings"

	response "task-api/internal/platform/http"
)

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type CreateTaskEndpoint struct {
	store Store
}

func NewCreateTaskEndpoint(store Store) *CreateTaskEndpoint {
	return &CreateTaskEndpoint{store: store}
}

// @Summary Create task
// @Tags tasks
// @Accept json
// @Produce json
// @Param input body CreateTaskRequest true "Create Payload"
// @Success 201 {object} Task
// @Failure 400 {object} response.ErrorResponse
// @Router /tasks [post]
func (e *CreateTaskEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		response.Error(w, http.StatusBadRequest, "Task title is required and cannot be empty")
		return
	}

	task, err := e.store.Create(req.Title)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	response.JSON(w, http.StatusCreated, task)
}
