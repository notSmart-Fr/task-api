package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	supabase "github.com/nedpals/supabase-go"

	// Anonymous import registers generated Swagger docs side-effects for OpenAPI tools
	_ "task-api/docs"

	// Feature modules (Vertical Slices)
	"task-api/internal/features/auth"
	featureDocs "task-api/internal/features/docs"
	"task-api/internal/features/health"
	"task-api/internal/features/system"
	"task-api/internal/features/tasks"
	"task-api/internal/features/users"

	// Platform infrastructure & middleware
	"task-api/internal/platform/middleware"
	"task-api/internal/platform/router"
)

// Global OpenAPI / Swagger Annotations used by 'swag init'
// @title Task & Auth API
// @version 1.0
// @description Production Go API with Supabase Auth & SQLite persistence.
// @Server http://localhost:8000 Local Development Server
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer " followed by your Supabase access_token
func main() {

	// -------------------------------------------------------------------------
	// 1. BOOTSTRAP ENVIRONMENT & STRUCTURED LOGGING
	// -------------------------------------------------------------------------
	// Load environment variables from .env into os.Getenv (silently ignores if .env is missing, e.g., in Docker)
	_ = godotenv.Load()

	// Set up global structured JSON logging to Stdout using Go's standard slog package
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// -------------------------------------------------------------------------
	// 2. INFRASTRUCTURE & EXTERNAL CLIENT INITIALIZATION
	// -------------------------------------------------------------------------
	// Read Supabase credentials for client creation
	sbURL := os.Getenv("SUPABASE_URL")
	sbPublishableKey := os.Getenv("SUPABASE_PUBLISHABLE_KEY")

	// Instantiate global Supabase SDK client for authentication features
	sbClient := supabase.CreateClient(sbURL, sbPublishableKey)

	// Instantiate Auth Middleware (lazy-loads & caches Supabase JWKS public keys on demand)
	authMw := middleware.NewAuthMiddleware(sbURL)

	// Instantiate HTTP Router Mux (Go 1.22+ standard HTTP router)
	mux := http.NewServeMux()

	// Read database path with fallback to default "tasks.db"
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "tasks.db"
	}

	// Initialize SQLite database connection pool and store for tasks
	taskStore, err := tasks.NewSQLiteStore(dbPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	// -------------------------------------------------------------------------
	// 3. FEATURE MODULE REGISTRATION & LIFECYCLE MANAGEMENT
	// -------------------------------------------------------------------------
	// Declare slice of feature modules implementing the router.Module interface
	modules := []router.Module{
		system.NewModule(),      // Base info / root endpoints
		health.NewModule(),      // Healthcheck endpoint
		featureDocs.NewModule(), // Scalar UI & openapi.json serving
		auth.NewModule(sbClient, authMw, sbURL, sbPublishableKey), // Auth slice (login, signup, logout)
		users.NewModule(authMw),                                   // Protected profile slice
		tasks.NewModule(taskStore),                                // CRUD task slice
	}

	// Iterate through the modules slice and attach all routes to 'mux'
	router.RegisterModules(mux, modules...)

	// Gracefully close resources held by feature modules (e.g., DB pools) when main exits
	defer func() {
		if err := router.CloseModules(modules...); err != nil {
			slog.Error("Failed to close feature modules", "error", err)
		}
	}()

	// -------------------------------------------------------------------------
	// 4. GLOBAL MIDDLEWARE PIPELINE
	// -------------------------------------------------------------------------
	// Wrap top-level router with global request logging middleware
	handler := middleware.Logging(mux)

	// -------------------------------------------------------------------------
	// 5. HTTP SERVER STARTUP
	// -------------------------------------------------------------------------
	// Determine port configuration with fallback to 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Boot HTTP server listener
	slog.Info("Server listening", "port", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
