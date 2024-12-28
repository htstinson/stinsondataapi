package main

import (
	"api/internal/auth"
	"api/internal/handler"
	"api/internal/middleware"
	"api/pkg/database"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecretString(secretName string, region string) (string, error) {

	var SecretValue string

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return SecretValue, err
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return SecretValue, err
	}

	var secretValue string = *result.SecretString

	return secretValue, err

}

func main() {
	// Create logger
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)

	rdsSecrets, err := GetSecretString("API", "us-west-2")
	if err != nil {
		logger.Println("RDS secrets", err.Error())
		return
	}

	fmt.Println(rdsSecrets)

	rdsPassword, err := GetSecretString("RDS/apidb", "us-west-2")
	if err != nil {
		logger.Println("RDS password", err.Error())
		return
	}

	fmt.Println(rdsPassword)

	if true {
		return
	}

	// Initialize database
	logger.Println("initializing database")
	db, err := database.New(database.Config{
		Host:     "", //rdsHost,
		Port:     5432,
		User:     "", //rdsUser,
		Password: "", //rdsPassword,
		DBName:   "", //rdsDbName,
		SSLMode:  "require",
	})
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize auth
	authConfig := auth.Config{
		SecretKey:     os.Getenv("JWT_SECRET_KEY"), // Use environment variable
		TokenDuration: 24 * time.Hour,              // Token valid for 24 hours
	}

	if authConfig.SecretKey == "" {
		authConfig.SecretKey = "your-secret-key-for-development" // Default for development
	}

	jwtAuth := auth.New(authConfig)

	// Create router and handler
	router := mux.NewRouter()
	//h := handler.NewHandler(db)

	// Create handler with auth
	h := handler.NewHandler(
		db,
		*jwtAuth,
		logger,
	)

	// Setup routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes
	api.HandleFunc("/health", h.HealthCheck).Methods("GET")
	api.HandleFunc("/register", h.Register).Methods("POST")
	api.HandleFunc("/login", h.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/", h.HealthCheck).Methods("GET")

	// Protected rounts
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(jwtAuth.Middleware)

	protected.HandleFunc("/items", h.CreateItem).Methods("POST", "OPTIONS")
	protected.HandleFunc("/items/{id}", h.UpdateItem).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/items/{id}", h.GetItem).Methods("GET")
	protected.HandleFunc("/items", h.ListItems).Methods("GET", "OPTIONS")

	protected.HandleFunc("/items/{id}", h.DeleteItem).Methods("DELETE")

	protected.HandleFunc("/accounts", h.CreateAccount).Methods("POST", "OPTIONS")
	protected.HandleFunc("/accounts/{id}", h.UpdateAccount).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/accounts", h.ListAccounts).Methods("GET", "OPTIONS")

	protected.HandleFunc("/users", h.CreateUser).Methods("POST", "OPTIONS")
	protected.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT", "OPTIONS")

	protected.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
	protected.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	protected.HandleFunc("/users", h.ListUsers).Methods("GET", "OPTIONS")

	// Add middleware
	api.Use(middleware.Logger(logger))
	api.Use(middleware.RequestID)
	api.Use(middleware.SecurityHeaders)
	api.Use(middleware.CORS)

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
		logger.Printf("Server starting.")

		err := srv.ListenAndServeTLS("certs/combined_certificate.crt", "certs/private.key")
		if err == http.ErrServerClosed {
			logger.Printf("Failed to start server (tls): %v", err)
		} else {
			err = srv.ListenAndServe()
			if err != nil {
				logger.Printf("Failed to start server: %v", err)
			}
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
