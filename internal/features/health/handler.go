package health

import (
	"net/http"

	response "task-api/internal/platform/http"
)

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
}

func (m *Module) Close() error {
	return nil
}
