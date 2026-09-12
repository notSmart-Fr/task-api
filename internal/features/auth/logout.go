package auth

import (
	"net/http"
	"strings"

	response "task-api/internal/platform/http"

	supabase "github.com/nedpals/supabase-go"
)

type LogoutEndpoint struct {
	client *supabase.Client
}

func NewLogoutEndpoint(client *supabase.Client) *LogoutEndpoint {
	return &LogoutEndpoint{client: client}
}

// @Summary Log out user
// @Tags auth
// @Security Bearer
// @Success 204 "User successfully logged out"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Logout failed"
// @Router /auth/logout [post]
func (e *LogoutEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if err := e.client.Auth.SignOut(r.Context(), token); err != nil {
		response.Error(w, http.StatusInternalServerError, "Logout failed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
