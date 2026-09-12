package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	response "task-api/internal/platform/http"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpEndpoint struct {
	supabaseURL string
	apiKey      string
	httpClient  *http.Client
}

func NewSignUpEndpoint(supabaseURL, apiKey string, httpClient *http.Client) *SignUpEndpoint {
	return &SignUpEndpoint{
		supabaseURL: strings.TrimRight(supabaseURL, "/"),
		apiKey:      apiKey,
		httpClient:  httpClient,
	}
}

// @Summary Sign up user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.SignUpRequest true "User credentials for registration"
// @Success 201 {object} map[string]any "User successfully created"
// @Failure 400 {object} response.ErrorResponse "Invalid input data or registration failure"
// @Router /auth/signup [post]
func (e *SignUpEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	body, err := json.Marshal(req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create signup request")
		return
	}

	request, err := http.NewRequestWithContext(r.Context(), http.MethodPost, e.supabaseURL+"/auth/v1/signup", bytes.NewReader(body))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create signup request")
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("apikey", e.apiKey)

	client := e.httpClient
	if client == nil {
		client = http.DefaultClient
	}
	upstreamResponse, err := client.Do(request)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "Unable to reach authentication service")
		return
	}
	defer upstreamResponse.Body.Close()

	var signUpResponse json.RawMessage
	if err := json.NewDecoder(upstreamResponse.Body).Decode(&signUpResponse); err != nil {
		response.Error(w, http.StatusBadGateway, "Invalid response from authentication service")
		return
	}
	if upstreamResponse.StatusCode < http.StatusOK || upstreamResponse.StatusCode >= http.StatusMultipleChoices {
		response.Error(w, upstreamResponse.StatusCode, fmt.Sprintf("Supabase signup failed: %s", signUpResponse))
		return
	}

	response.JSON(w, http.StatusCreated, signUpResponse)
}
