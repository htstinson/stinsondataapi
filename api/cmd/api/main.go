package main

import (
	"fmt"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/auth"
	"github.com/htstinson/stinsondataapi/api/internal/handler"
	"github.com/htstinson/stinsondataapi/api/internal/middleware"
	"github.com/htstinson/stinsondataapi/api/internal/model"
	"github.com/htstinson/stinsondataapi/api/pkg/database"
	"github.com/htstinson/stinsondataapi/api/salesforce"

	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"log/syslog"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	// Register the correct MIME types for JavaScript modules
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".mjs", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
}

func main() {
	// Create logger
	logger, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "webserver")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log.SetOutput(logger)

	fmt.Println("initialize salesforce")
	sf, err := salesforce.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("[%v] Initializing database", time.Now().Format(time.RFC3339))
	var RDSLogin = &model.RDSLogin{}
	rdsLogin, err := common.GetSecretString("RDS/apidb", "us-west-2")
	if err != nil {
		fmt.Println("RDS Login", err.Error())
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
		fmt.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	fmt.Println("Connected to Database")

	// Initialize auth
	authConfig := auth.Config{
		SecretKey:     os.Getenv("JWT_SECRET_KEY"), // Use environment variable
		TokenDuration: 24 * time.Hour,              // Token valid for 24 hours
	}

	if authConfig.SecretKey == "" {
		authConfig.SecretKey = "your-secret-key-for-development" // Default for development
	}

	jwtAuth := auth.New(authConfig)
	// Create handler with auth and SFauth
	h := handler.NewHandler(db, *jwtAuth, &log.Logger{})

	// Create router and handler
	router := mux.NewRouter()
	// debug region start

	//test := router.PathPrefix("/test").Subrouter()           //
	//test.HandleFunc("/health", h.HealthCheck).Methods("GET") //
	//test.HandleFunc("/users", h.ListUsers).Methods("GET", "OPTIONS")

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

	protected.HandleFunc("/accounts", sf.Handler.CreateAccount).Methods("POST", "OPTIONS")
	protected.HandleFunc("/accounts/{id}", sf.Handler.UpdateAccount).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/accounts", sf.Handler.ListAccounts).Methods("GET", "OPTIONS")

	protected.HandleFunc("/contacts", sf.Handler.ListContacts).Methods("GET", "OPTIONS")
	protected.HandleFunc("/contacts/{accountid}", sf.Handler.ListContacts).Methods("GET", "OPTIONS")
	protected.HandleFunc("/contact/{contactid}", sf.Handler.GetContactById).Methods("GET", "OPTIONS")

	protected.HandleFunc("/users", h.CreateUser).Methods("POST", "OPTIONS")
	protected.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
	protected.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	protected.HandleFunc("/users", h.ListUsers).Methods("GET", "OPTIONS")

	// Add middleware
	api.Use(middleware.Logger(&log.Logger{}))
	api.Use(middleware.RequestID)
	api.Use(middleware.SecurityHeaders)
	api.Use(middleware.CORS)

	//static assets
	distPath := "/home/ec2-user/go/src/stinsondata-tools-reactapp/dist"
	fmt.Printf("Serving files from: %s", distPath)

	// Handle all static assets including the index.js file
	router.PathPrefix("/assets/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Asset request for: %s", r.URL.Path)

		// Remove the leading /assets/ to get the file path
		filePath := filepath.Join(distPath, r.URL.Path)
		fmt.Printf("Looking for file at: %s", filePath)

		// Set appropriate headers based on file extension
		switch ext := path.Ext(r.URL.Path); ext {
		case ".js", ".mjs", ".jsx":
			w.Header().Set("Content-Type", "application/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		}

		// Serve the file
		http.ServeFile(w, r, filePath)
	})

	// Handle root and all other routes with index.html
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Serving index.html for path: %s", r.URL.Path)
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
	})

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
		fmt.Printf("Server starting.")

		err := srv.ListenAndServeTLS("../../certs/certificate.crt", "../../certs/private.key")
		if err == http.ErrServerClosed {
			fmt.Printf("Failed to start server (tls): %v", err)
		} else {
			fmt.Println(err.Error())
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	fmt.Println("Server stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v", err)
		return
	}

	fmt.Println("Server stopped")
}
