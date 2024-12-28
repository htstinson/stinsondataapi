package main

import (
	"api/internal/auth"
	"api/internal/handler"
	"api/internal/middleware"
	"api/internal/model"
	"api/pkg/database"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Basic HTML template that loads your React app
const indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Kendo React App</title>
    <!-- KendoReact styles -->
    <link href="https://kendo.cdn.telerik.com/themes/6.4.0/default/default-main.css" rel="stylesheet" />
    <!-- Your bundled React app -->
    <script src="/static/js/bundle.js" defer></script>
</head>
<body>
    <div id="root">
	Hi
	</div>
</body>
</html>
`

func init() {
	// Register the correct MIME types for JavaScript modules
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".mjs", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
}

func main() {
	// Create logger
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)

	logger.Println("initializing database")
	var RDSLogin = &model.RDSLogin{}
	rdsLogin, err := GetSecretString("RDS/apidb", "us-west-2")
	if err != nil {
		logger.Println("RDS Login", err.Error())
		return
	}
	json.Unmarshal(rdsLogin, RDSLogin)
	// Initialize database
	db, err := database.New(database.Config{
		Host:     RDSLogin.Host,
		Port:     RDSLogin.Port,
		User:     RDSLogin.Username,
		Password: RDSLogin.Password,
		DBName:   "apidb",
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
	// Create handler with auth
	h := handler.NewHandler(
		db,
		*jwtAuth,
		logger,
	)

	// Create router and handler
	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("../../../../stinsondata-tools-reactapp/dist"))

	// Wrap the file server with MIME type headers
	fsWithMime := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Requested path: %s\n", r.URL.Path)

		w.Header().Set("Cache-Control", "no-cache")

		switch ext := path.Ext(r.URL.Path); ext {
		case ".js", ".mjs", ".jsx":
			w.Header().Set("Content-Type", "application/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		}

		fs.ServeHTTP(w, r)
	})

	router.PathPrefix("/assets/").Handler( // Change from /static/ to /assets/
		http.StripPrefix("/assets/", fsWithMime),
	)

	// Index.html handler remains structurally the same
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "../../../../stinsondata-tools-reactapp/dist/index.html") // Update path
	})

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
			MinVersion: tls.VersionTLS13,
		},
	}

	// Start server
	go func() {
		logger.Printf("Server starting.")

		err := srv.ListenAndServeTLS("../../certs/certificate.crt", "../../certs/private.key")
		if err == http.ErrServerClosed {
			logger.Printf("Failed to start server (tls): %v", err)
		} else {
			logger.Println(err.Error())
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

func GetSecretString(secretName string, region string) ([]byte, error) {

	var SecretValue []byte

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

	SecretValue = []byte(*result.SecretString)

	return SecretValue, err

}
