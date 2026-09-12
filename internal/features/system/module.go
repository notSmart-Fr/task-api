package system

import (
	"net/http"

	response "task-api/internal/platform/http"
)

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]any{
			"name":      "Task & Auth API",
			"version":   "1.0",
			"endpoints": []string{"/tasks", "/auth/signup", "/auth/login", "/protected/profile", "/docs"},
		})
	})
	mux.HandleFunc("GET /public/info", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"message": "Welcome stranger! This info is public."})
	})
}

func (m *Module) Close() error {
	return nil
}
