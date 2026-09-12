package auth

import (
	"net/http"

	"task-api/internal/platform/middleware"

	supabase "github.com/nedpals/supabase-go"
)

type AuthModule struct {
	client      *supabase.Client
	authMw      *middleware.AuthMiddleware
	supabaseURL string
	apiKey      string
}

func NewModule(client *supabase.Client, authMw *middleware.AuthMiddleware, supabaseURL, apiKey string) *AuthModule {
	return &AuthModule{
		client:      client,
		authMw:      authMw,
		supabaseURL: supabaseURL,
		apiKey:      apiKey,
	}
}

func (m *AuthModule) RegisterRoutes(mux *http.ServeMux) {
	signUpEp := NewSignUpEndpoint(m.supabaseURL, m.apiKey, m.client.HTTPClient)
	loginEp := NewLoginEndpoint(m.client)
	logoutEp := NewLogoutEndpoint(m.client)

	mux.HandleFunc("POST /auth/signup", signUpEp.Handle)
	mux.HandleFunc("POST /auth/login", loginEp.Handle)
	mux.HandleFunc("POST /auth/logout", m.authMw.RequireAuth(logoutEp.Handle))
}

func (m *AuthModule) Close() error {
	return nil
}
