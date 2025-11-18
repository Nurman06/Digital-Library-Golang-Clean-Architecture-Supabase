package main

import (
	"net/http"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/adapter/handler"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/config"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/server"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/supabase"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase"
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

	// Initialize repositories
	bookRepo := repository.NewPostgresBookRepository(supabaseClient.DB)
	userRepo := repository.NewPostgresUserRepository(supabaseClient.DB)
	borrowRecordRepo := repository.NewPostgresBorrowRecordRepository(supabaseClient.DB)
	bookCopyRepo := repository.NewPostgresBookCopyRepository(supabaseClient.DB)
	// reservationRepo := repository.NewPostgresReservationRepository(supabaseClient.DB)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, borrowRecordRepo)
	bookUseCase := usecase.NewBookUseCase(bookRepo, bookCopyRepo)
	borrowingUseCase := usecase.NewBorrowingUseCase(borrowRecordRepo, bookCopyRepo, userRepo, bookRepo)
	// searchUseCase := usecase.NewSearchUseCase(bookRepo, bookCopyRepo, borrowRecordRepo)
	// availabilityUseCase := usecase.NewAvailabilityUseCase(bookRepo, bookCopyRepo, borrowRecordRepo, reservationRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase, userUseCase, appLogger)
	bookHandler := handler.NewBookHandler(bookUseCase, authUseCase, appLogger)
	borrowingHandler := handler.NewBorrowingHandler(borrowingUseCase, bookUseCase, userUseCase, appLogger)
	userHandler := handler.NewUserHandler(userUseCase, authUseCase, appLogger)
	availabilityHandler := handler.NewAvailabilityHandler(bookUseCase, appLogger)

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
	setupRoutes(router, authHandler, bookHandler, borrowingHandler, userHandler, availabilityHandler, appLogger)

	// Start server
	appLogger.Infof("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := httpServer.Start(); err != nil {
		appLogger.Fatalf("Server failed: %v", err)
	}
}

func setupRoutes(
	router *mux.Router,
	authHandler *handler.AuthHandler,
	bookHandler *handler.BookHandler,
	borrowingHandler *handler.BorrowingHandler,
	userHandler *handler.UserHandler,
	availabilityHandler *handler.AvailabilityHandler,
	appLogger *logger.Logger,
) {
	// Apply global middleware
	router.Use(handler.CORSMiddleware)
	router.Use(handler.LoggingMiddleware(appLogger))
	router.Use(handler.RecoveryMiddleware(appLogger))

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"digital-library-api"}`))
	}).Methods("GET")

	// API v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Root endpoint
	apiV1.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Digital Library API v1","status":"ready","version":"1.0.0"}`))
	}).Methods("GET")

	// Auth routes (public)
	authRoutes := apiV1.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")
	
	// Protected auth routes
	authProtected := authRoutes.PathPrefix("").Subrouter()
	authProtected.Use(handler.AuthMiddleware)
	authProtected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")

	// Book routes
	bookRoutes := apiV1.PathPrefix("/books").Subrouter()
	
	// Public book routes
	bookRoutes.HandleFunc("", bookHandler.ListBooks).Methods("GET")
	bookRoutes.HandleFunc("/search", bookHandler.SearchBooks).Methods("GET")
	bookRoutes.HandleFunc("/{id}", bookHandler.GetBook).Methods("GET")
	
	// Book availability routes (public)
	bookRoutes.HandleFunc("/{id}/availability", availabilityHandler.GetBookAvailability).Methods("GET")
	bookRoutes.HandleFunc("/{id}/copies", availabilityHandler.GetBookCopies).Methods("GET")
	bookRoutes.HandleFunc("/{id}/copies/available", availabilityHandler.GetAvailableBookCopies).Methods("GET")
	
	// Protected book routes (Admin/Librarian only)
	bookProtected := bookRoutes.PathPrefix("").Subrouter()
	bookProtected.Use(handler.AuthMiddleware)
	bookProtected.Use(handler.RequireRole("admin", "librarian"))
	bookProtected.HandleFunc("", bookHandler.CreateBook).Methods("POST")
	bookProtected.HandleFunc("/{id}", bookHandler.UpdateBook).Methods("PUT")
	bookProtected.HandleFunc("/{id}", bookHandler.DeleteBook).Methods("DELETE")

	// Borrowing routes (all protected)
	borrowingRoutes := apiV1.PathPrefix("/borrowing").Subrouter()
	borrowingRoutes.Use(handler.AuthMiddleware)
	borrowingRoutes.HandleFunc("/checkout", borrowingHandler.CheckoutBook).Methods("POST")
	borrowingRoutes.HandleFunc("/return", borrowingHandler.ReturnBook).Methods("POST")
	borrowingRoutes.HandleFunc("/renew", borrowingHandler.RenewBook).Methods("POST")
	borrowingRoutes.HandleFunc("/history", borrowingHandler.GetBorrowingHistory).Methods("GET")
	borrowingRoutes.HandleFunc("/active", borrowingHandler.GetActiveBorrows).Methods("GET")
	borrowingRoutes.HandleFunc("/overdue", borrowingHandler.GetOverdueBorrows).Methods("GET")
	borrowingRoutes.HandleFunc("/{id}", borrowingHandler.GetBorrowRecord).Methods("GET")

	// User routes
	userRoutes := apiV1.PathPrefix("/users").Subrouter()
	userRoutes.Use(handler.AuthMiddleware)
	
	// User profile routes (any authenticated user)
	userRoutes.HandleFunc("/profile", userHandler.UpdateUserProfile).Methods("PUT")
	
	// Admin-only user routes
	userAdminRoutes := userRoutes.PathPrefix("").Subrouter()
	userAdminRoutes.Use(handler.RequireRole("admin"))
	userAdminRoutes.HandleFunc("", userHandler.ListUsers).Methods("GET")
	userAdminRoutes.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	userAdminRoutes.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	userAdminRoutes.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")
	userAdminRoutes.HandleFunc("/{id}/suspend", userHandler.SuspendUser).Methods("POST")
	userAdminRoutes.HandleFunc("/{id}/activate", userHandler.ActivateUser).Methods("POST")

	appLogger.Info("Routes configured successfully")
}
