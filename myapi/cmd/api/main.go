package main

import (
	"context"
	"crypto/tls"
	"log"
	"myapi/internal/handler"
	"myapi/internal/middleware"
	"myapi/pkg/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Create logger
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)

	// Initialize database
	db, err := database.New(database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "apidb",
		SSLMode:  "disable",
	})
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create router and handler
	router := mux.NewRouter()
	h := handler.NewHandler(db)

	// Setup routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/health", h.HealthCheck).Methods("GET")
	api.HandleFunc("/items", h.CreateItem).Methods("POST")
	api.HandleFunc("/items/{id}", h.GetItem).Methods("GET")
	api.HandleFunc("/items", h.ListItems).Methods("GET")

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.RequestID)
	router.Use(middleware.SecurityHeaders)
	router.Use(middleware.CORS)

	// Create server with local certificates
	srv := &http.Server{
		Addr:    ":443",
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	// Start server
	go func() {
		logger.Printf("Server starting on https://api.local.dev")
		if err := srv.ListenAndServeTLS("certs/local.crt", "certs/local.key"); err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Println("Server stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server stopped")
}
