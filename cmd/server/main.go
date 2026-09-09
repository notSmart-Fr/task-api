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

// @title Task API
// @version 1.0
// @description Scalable Task CRUD API built with Go and SQLite.
// @host localhost:8000
// @BasePath /
func main() {
	// Initialize structured JSON logging to stdout
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create a new Go 1.22+ standard library HTTP multiplexer
	mux := http.NewServeMux()

	// Root Endpoint: returns basic API metadata
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

	// Initialize SQLite Database Store (creates tasks.db and table if missing)
	taskStore, err := tasks.NewSQLiteStore("tasks.db")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Pass taskStore directly because *tasks.SQLiteStore implements tasks.Store
	tasks.RegisterHandlers(mux, taskStore)

	// Wrap root multiplexer with global logging middleware
	handler := middleware.Logging(mux)

	port := ":8000"
	slog.Info("Server listening", "port", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
