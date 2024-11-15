package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"myapi/internal/handler"
	"myapi/internal/middleware"
	"myapi/pkg/database"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
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

	// Create router
	router := mux.NewRouter()

	// Create handler
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

	// Create server
	srv := &http.Server{
		Addr:         ":443",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Configure TLS
	m := &autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("api.local.dev"),
	}
	srv.TLSConfig = m.TLSConfig()

	// Start HTTP server (redirect to HTTPS)
	go http.ListenAndServe(":80", m.HTTPHandler(nil))

	// Start server
	go func() {
		logger.Printf("Server starting on https://api.local.dev")
		if err := srv.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
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
