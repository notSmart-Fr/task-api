package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	response "task-api/internal/platform/http"
)

type userContextKey string

const (
	UserIDKey    userContextKey = "user_id"
	UserEmailKey userContextKey = "user_email"
)

type AuthMiddleware struct {
	supabaseURL string
	apiKey      string
}

func NewAuthMiddleware(supabaseURL string) *AuthMiddleware {
	apiKey := os.Getenv("SUPABASE_PUBLISHABLE_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("SUPABASE_KEY")
	}

	return &AuthMiddleware{
		supabaseURL: supabaseURL,
		apiKey:      apiKey,
	}
}

type authUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "Access token required")
			return
		}

		tokenString, ok := bearerToken(authHeader)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		user, err := m.getUser(r, tokenString)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
		if user.ID == "" {
			response.Error(w, http.StatusUnauthorized, "Missing user ID in token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
		ctx = context.WithValue(ctx, UserEmailKey, user.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) getUser(request *http.Request, token string) (authUser, error) {
	url := strings.TrimRight(m.supabaseURL, "/") + "/auth/v1/user"
	upstreamRequest, err := http.NewRequestWithContext(request.Context(), http.MethodGet, url, nil)
	if err != nil {
		return authUser{}, fmt.Errorf("create Supabase user request: %w", err)
	}
	upstreamRequest.Header.Set("Authorization", "Bearer "+token)
	upstreamRequest.Header.Set("apikey", m.apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	upstreamResponse, err := client.Do(upstreamRequest)
	if err != nil {
		return authUser{}, fmt.Errorf("request Supabase user: %w", err)
	}
	defer upstreamResponse.Body.Close()
	if upstreamResponse.StatusCode != http.StatusOK {
		return authUser{}, fmt.Errorf("Supabase user request returned status %d", upstreamResponse.StatusCode)
	}

	var user authUser
	if err := json.NewDecoder(upstreamResponse.Body).Decode(&user); err != nil {
		return authUser{}, fmt.Errorf("decode Supabase user: %w", err)
	}
	return user, nil
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	return parts[1], true
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}

func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}
