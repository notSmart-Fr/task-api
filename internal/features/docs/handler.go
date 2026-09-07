package docs

import (
	"fmt"
	"net/http"
	"os"
)

func RegisterHandlers(mux *http.ServeMux) {
	// Serve raw OpenAPI JSON spec
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("openapi.json")
		if err != nil {
			http.Error(w, "Unable to load API specification", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	// Serve Scalar API Reference UI at /docs
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
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, html)
	})
}
