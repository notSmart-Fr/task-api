package users

import (
	"net/http"

	"task-api/internal/platform/middleware"
)

type Module struct {
	authMw *middleware.AuthMiddleware
}

func NewModule(authMw *middleware.AuthMiddleware) *Module {
	return &Module{authMw: authMw}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	profileEp := NewProfileEndpoint()
	mux.HandleFunc("GET /protected/profile", m.authMw.RequireAuth(profileEp.Handle))
}

func (m *Module) Close() error {
	return nil
}
