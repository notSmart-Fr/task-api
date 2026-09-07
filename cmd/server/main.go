package main

import (
	"fmt"
	"net/http"

	"task-api/internal/features/health"
	response "task-api/internal/platform/http"
)

func main() {
	mux := http.NewServeMux()

	// Stage 1: Root Endpoint returning API description JSON
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]any{
			"name":      "Task API",
			"version":   "1.0",
			"endpoints": []string{"/tasks"},
		})
	})

	// Stage 1: Register Health Slice
	health.RegisterHandlers(mux)

	port := ":8000"
	fmt.Printf("Server starting on http://localhost%s...\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
