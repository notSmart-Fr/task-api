package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	supabase "github.com/nedpals/supabase-go"

	_ "task-api/docs"
	"task-api/internal/features/auth"
	featureDocs "task-api/internal/features/docs"
	"task-api/internal/features/health"
	"task-api/internal/features/system"
	"task-api/internal/features/tasks"
	"task-api/internal/features/users"
	"task-api/internal/platform/middleware"
	"task-api/internal/platform/router"
)

// @title Task & Auth API
// @version 1.0
// @description Production Go API with Supabase Auth & SQLite persistence.
// @Server http://localhost:8000 Local Development Server
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer " followed by your Supabase access_token
func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	sbURL := os.Getenv("SUPABASE_URL")
	sbPublishableKey := os.Getenv("SUPABASE_PUBLISHABLE_KEY")

	sbClient := supabase.CreateClient(sbURL, sbPublishableKey)

	// Direct initialization — no error check needed because public keys
	// are fetched and cached cleanly on the first authenticated request.
	authMw := middleware.NewAuthMiddleware(sbURL)

	mux := http.NewServeMux()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "tasks.db"
	}
	taskStore, err := tasks.NewSQLiteStore(dbPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	modules := []router.Module{
		system.NewModule(),
		health.NewModule(),
		featureDocs.NewModule(),
		auth.NewModule(sbClient, authMw, sbURL, sbPublishableKey),
		users.NewModule(authMw),
		tasks.NewModule(taskStore),
	}
	router.RegisterModules(mux, modules...)
	defer func() {
		if err := router.CloseModules(modules...); err != nil {
			slog.Error("Failed to close feature modules", "error", err)
		}
	}()

	handler := middleware.Logging(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	slog.Info("Server listening", "port", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
