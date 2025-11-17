package main

import (
	"net/http"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/config"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/server"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/supabase"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize logger
	appLogger := logger.New()
	appLogger.Info("Starting Digital Library API...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatalf("Failed to load configuration: %v", err)
	}
	appLogger.Infof("Configuration loaded successfully (Environment: %s)", cfg.Server.Env)

	// Initialize Supabase client
	supabaseClient, err := supabase.NewClient(supabase.Config{
		URL:             cfg.Supabase.URL,
		AnonKey:         cfg.Supabase.AnonKey,
		ServiceKey:      cfg.Supabase.ServiceKey,
		DatabaseDSN:     cfg.GetDatabaseDSN(),
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute,
	})
	if err != nil {
		appLogger.Fatalf("Failed to initialize Supabase client: %v", err)
	}
	defer supabaseClient.Close()
	appLogger.Info("Supabase client initialized successfully")

	// Test database connection
	if err := supabaseClient.Health(); err != nil {
		appLogger.Fatalf("Database health check failed: %v", err)
	}
	appLogger.Info("Database connection established")

	// Initialize HTTP server
	httpServer := server.NewServer(server.Config{
		Host:         cfg.Server.Host,
		Port:         cfg.Server.Port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// Get router
	router := httpServer.GetRouter()

	// Setup routes
	setupRoutes(router, supabaseClient, cfg, appLogger)

	// Start server
	appLogger.Infof("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := httpServer.Start(); err != nil {
		appLogger.Fatalf("Server failed: %v", err)
	}
}

func setupRoutes(router *mux.Router, _ *supabase.Client, _ *config.Config, appLogger *logger.Logger) {
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"digital-library-api"}`))
	}).Methods("GET")

	// API v1 routes will be added here
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Placeholder for future routes
	apiV1.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Digital Library API v1","status":"ready"}`))
	}).Methods("GET")

	appLogger.Info("Routes configured successfully")
}
