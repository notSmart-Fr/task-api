package main

import (
	"log/slog"
	"net/http"
	"os"

	"task-api/internal/features/docs"
	"task-api/internal/features/health"
	"task-api/internal/features/tasks"
	response "task-api/internal/platform/http"
	"task-api/internal/platform/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]any{
			"name":      "Task API",
			"version":   "1.0",
			"endpoints": []string{"/tasks", "/health", "/docs"},
		})
	})

	// Register Feature Slices
	health.RegisterHandlers(mux)
	docs.RegisterHandlers(mux)

	// Initialize SQLite Store instead of MemoryStore
	taskStore, err := tasks.NewSQLiteStore("tasks.db")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	tasks.RegisterHandlers(mux, taskStore)

	handler := middleware.Logging(mux)

	port := ":8000"
	slog.Info("Server listening", "port", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
