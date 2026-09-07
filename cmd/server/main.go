package main

import (
	"log/slog"
	"net/http"
	"os"

	"task-api/internal/features/health"
	"task-api/internal/features/tasks"
	response "task-api/internal/platform/http"
	"task-api/internal/platform/middleware"
)

func main() {
	// Configure JSON logger output
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()

	// Root Route
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

	// Apply Global Middleware
	handler := middleware.Logging(mux)

	port := ":8000"
	slog.Info("Server listening", "port", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
