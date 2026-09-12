package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	response "task-api/internal/platform/http"
)

type DeleteTaskEndpoint struct {
	store Store
}

func NewDeleteTaskEndpoint(store Store) *DeleteTaskEndpoint {
	return &DeleteTaskEndpoint{store: store}
}

// @Summary Delete task
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204 "Task deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid task ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to delete task"
// @Failure 404 {object} response.ErrorResponse "Task not found"
// @Router /tasks/{id} [delete]
func (e *DeleteTaskEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	if err := e.store.Delete(id); err == ErrNotFound {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Task %d not found", id))
		return
	} else if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
