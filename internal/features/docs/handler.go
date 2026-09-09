package docs

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func RegisterHandlers(mux *http.ServeMux) {
	// Reads swagger.json from the root docs/ folder
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			slog.Error("Failed to read swagger.json from root docs folder", "error", err)
			http.Error(w, `{"error":"Unable to load API specification"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	// Serves Scalar UI at /docs
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		html := `<!doctype html>
<html>
  <head>
    <title>Task API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script id="api-reference" data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, html)
	})
}
