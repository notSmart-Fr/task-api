package users

import (
	"net/http"

	response "task-api/internal/platform/http"
	"task-api/internal/platform/middleware"
)

type ProfileEndpoint struct{}

func NewProfileEndpoint() *ProfileEndpoint {
	return &ProfileEndpoint{}
}

// @Summary Get current user profile
// @Tags protected
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]any "Current user details"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or missing token"
// @Router /protected/profile [get]
func (e *ProfileEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	email, _ := middleware.GetUserEmailFromContext(r.Context())

	response.JSON(w, http.StatusOK, map[string]any{
		"id":    userID,
		"email": email,
	})
}
