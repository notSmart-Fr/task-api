package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	response "task-api/internal/platform/http"

	supabase "github.com/nedpals/supabase-go"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginEndpoint struct {
	client *supabase.Client
}

func NewLoginEndpoint(client *supabase.Client) *LoginEndpoint {
	return &LoginEndpoint{client: client}
}

// @Summary Log in user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Credentials"
// @Success 200 {object} map[string]any
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
func (e *LoginEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	session, err := e.client.Auth.SignIn(r.Context(), supabase.UserCredentials{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid login credentials")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"access_token":  session.AccessToken,
		"refresh_token": session.RefreshToken,
		"user":          session.User,
	})
}
