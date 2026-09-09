package tasks

import (
	"net/http"
	"strconv"

	response "task-api/internal/platform/http"
)

type ListTasksRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginatedTasksResponse struct {
	Data       []Task `json:"data"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalItems int    `json:"total_items"`
	TotalPages int    `json:"total_pages"`
}

type ListTasksEndpoint struct {
	store Store
}

func NewListTasksEndpoint(store Store) *ListTasksEndpoint {
	return &ListTasksEndpoint{store: store}
}

// @Summary List all tasks with pagination
// @Tags tasks
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedTasksResponse
// @Router /tasks [get]
func (e *ListTasksEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	// Query parameters default handling
	page := 1
	limit := 10

	query := r.URL.Query()
	if pageStr := query.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	tasksList, total, err := e.store.GetAll(limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	res := PaginatedTasksResponse{
		Data:       tasksList,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}

	response.JSON(w, http.StatusOK, res)
}
