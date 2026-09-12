package docs

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			slog.Error("Failed to read swagger.json", "error", err)
			http.Error(w, `{"error":"Unable to load API specification"}`, http.StatusInternalServerError)
			return
		}
		data, err = normalizeOpenAPIDocument(data)
		if err != nil {
			slog.Error("Failed to normalize OpenAPI specification", "error", err)
			http.Error(w, `{"error":"Unable to load API specification"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		html := `<!doctype html>
<html>
  <head>
    <title>Task & Auth API Docs</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script 
      id="api-reference" 
      data-url="/openapi.json"
      data-configuration='{"servers":[{"url":"/"}]}'>
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, html)
	})
}

func (m *Module) Close() error {
	return nil
}

func normalizeOpenAPIDocument(data []byte) ([]byte, error) {
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, err
	}

	paths, _ := document["paths"].(map[string]any)
	for _, pathItem := range paths {
		operations, _ := pathItem.(map[string]any)
		for _, operation := range operations {
			normalizeRequestBody(operation)
		}
	}
	normalizeBearerSecurityScheme(document)

	return json.Marshal(document)
}

func normalizeBearerSecurityScheme(document map[string]any) {
	components, _ := document["components"].(map[string]any)
	securitySchemes, _ := components["securitySchemes"].(map[string]any)
	if _, exists := securitySchemes["Bearer"]; exists {
		securitySchemes["Bearer"] = map[string]any{
			"type":         "http",
			"scheme":       "bearer",
			"bearerFormat": "JWT",
		}
	}
}

func normalizeRequestBody(operation any) {
	operationMap, _ := operation.(map[string]any)
	requestBody, _ := operationMap["requestBody"].(map[string]any)
	content, _ := requestBody["content"].(map[string]any)

	for _, mediaType := range content {
		mediaTypeMap, _ := mediaType.(map[string]any)
		schema, _ := mediaTypeMap["schema"].(map[string]any)
		oneOf, _ := schema["oneOf"].([]any)
		if len(oneOf) != 2 {
			continue
		}

		emptyObject, _ := oneOf[0].(map[string]any)
		reference, _ := oneOf[1].(map[string]any)
		ref, _ := reference["$ref"].(string)
		if len(emptyObject) == 1 && emptyObject["type"] == "object" && ref != "" {
			mediaTypeMap["schema"] = map[string]any{"$ref": ref}
		}
	}
}
