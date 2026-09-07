package health

import (
	"net/http"

	response "task-api/internal/platform/http"
)

// RegisterHandlers registers the GET /health route on the provided ServeMux.
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
}
