package main

import (
	"fmt"
	"net/http"

	"task-api/internal/features/health"
	"task-api/internal/features/tasks"
	response "task-api/internal/platform/http"
)

func main() {
	mux := http.NewServeMux()

	// Root Endpoint
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]any{
			"name":      "Task API",
			"version":   "1.0",
			"endpoints": []string{"/tasks"},
		})
	})

	// Register Features
	health.RegisterHandlers(mux)

	taskStore := tasks.NewMemoryStore()
	tasks.RegisterHandlers(mux, taskStore)

	port := ":8000"
	fmt.Printf("Server starting on http://localhost%s...\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
