package main

import (
	"fmt"
	"net"

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

	fmt.Printf("\n[%v] ---------------------------------------------------------------\n", time.Now().Format(time.RFC3339))

	log.SetOutput(os.Stdout)

	fmt.Printf("[%v] Initializing salesforce\n", time.Now().Format(time.RFC3339))

	sf, err := salesforce.New()
	if err != nil {
		fmt.Printf("[%v] Salesforce error: %s\n", time.Now().Format(time.RFC3339), err.Error())
		return
	}

	fmt.Printf("[%v] Initializing database\n", time.Now().Format(time.RFC3339))
	var RDSLogin = &model.RDSLogin{}
	rdsLogin, err := common.GetSecretString("RDS/apidb", "us-west-2")
	if err != nil {
		fmt.Printf("[%v] RDS error: %s\n", time.Now().Format(time.RFC3339), err.Error())
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
		fmt.Printf("[%v] Failed to connect to database: %s\n", time.Now().Format(time.RFC3339), err.Error())
		return
	}
	defer db.Close()

	fmt.Printf("[%v] Connected to database\n", time.Now().Format(time.RFC3339))

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

	// Setup routes
	api := router.PathPrefix("/api/v1").Subrouter()

	api.Use(ipLoggingMiddleware)

	// Public routes
	api.HandleFunc("/health", h.HealthCheck).Methods("GET")
	api.HandleFunc("/register", h.Register).Methods("POST")
	api.HandleFunc("/login", h.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/", h.HealthCheck).Methods("GET")

	// Protected rounts
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(ipLoggingMiddleware)
	protected.Use(jwtAuth.Middleware)

	protected.HandleFunc("/admin", h.ListBlocked).Methods("GET", "OPTIONS")

	protected.HandleFunc("/blocked", h.ListBlocked).Methods("GET", "OPTIONS")
	protected.HandleFunc("/blocked/{id}", h.UpdateBlocked).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/blocked", h.CreateBlocked).Methods("POST", "OPTIONS")
	protected.HandleFunc("/blocked/{id}", h.DeleteBlocked).Methods("DELETE")

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
	fmt.Printf("[%v] Serving files from: %s\n", time.Now().Format(time.RFC3339), distPath)

	// Handle all static assets including the index.js file
	router.PathPrefix("/assets/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%v] Asset request for: %s\n", time.Now().Format(time.RFC3339), r.URL.Path)

		// Remove the leading /assets/ to get the file path
		filePath := filepath.Join(distPath, r.URL.Path)
		fmt.Printf("[%v] Looking for file at: %s\n", time.Now().Format(time.RFC3339), filePath)

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
		//fmt.Printf("[%v] Serving index.html for path: %s\n", time.Now().Format(time.RFC3339), r.URL.Path)
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
		fmt.Printf("[%v] Server starting\n", time.Now().Format(time.RFC3339))

		err := srv.ListenAndServeTLS("../../certs/certificate.crt", "../../certs/private.key")
		if err == http.ErrServerClosed {
			fmt.Printf("[%v] Failed to start server (tls): %v\n", time.Now().Format(time.RFC3339), err.Error())
		} else {
			fmt.Printf("[%v] Error: %s\n", time.Now().Format(time.RFC3339), err.Error())
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	fmt.Printf("[%v] Server stopping...\n", time.Now().Format(time.RFC3339))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("[%v] Server forced to shutdown: %s\n", time.Now().Format(time.RFC3339), err.Error())
		return
	}

	fmt.Printf("[%v] Server stopped.\n", time.Now().Format(time.RFC3339))
}

// Get the direct TCP/IP connection address (Layer 3)
func getTCPAddr(r *http.Request) string {
	// RemoteAddr contains the actual TCP connection address (IP:port)
	// This is the most reliable source of the client's direct IP
	// but will be the proxy's IP if the client is behind a proxy
	addr := r.RemoteAddr

	// RemoteAddr includes both IP and port (e.g., 192.168.1.1:12345)
	// Extract just the IP part
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		// If there's an error splitting, just return the whole thing
		return addr
	}

	return ip
}

// Log middleware that captures the TCP address
func ipLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the TCP address
		ipAddr := getTCPAddr(r)

		// Log the connection information
		log.Printf("Layer 3 connection from: %s, Method: %s, Path: %s",
			ipAddr, r.Method, r.URL.Path)

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// Handler function that displays the IP information
func ipInfoHandler(w http.ResponseWriter, r *http.Request) {
	ipAddr := getTCPAddr(r)

	// Is this a private IP address?
	ip := net.ParseIP(ipAddr)
	isPrivate := isPrivateIP(ip)

	fmt.Fprintf(w, "Connection Information\n\n")
	fmt.Fprintf(w, "Your TCP/IP address: %s\n", ipAddr)
	fmt.Fprintf(w, "Is private address: %t\n", isPrivate)

	// For educational purposes, also show what the headers claim
	// (but we're not using these for our actual IP detection)
	fmt.Fprintf(w, "\nHTTP Headers (NOT TRUSTED):\n")
	fmt.Fprintf(w, "X-Forwarded-For: %s\n", r.Header.Get("X-Forwarded-For"))
	fmt.Fprintf(w, "X-Real-IP: %s\n", r.Header.Get("X-Real-IP"))
}

// Check if an IP is a private address
func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	// Private IPv4 ranges
	privateRanges := []struct {
		start net.IP
		end   net.IP
	}{
		{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
		{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
		{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
		{net.ParseIP("127.0.0.0"), net.ParseIP("127.255.255.255")},
	}

	for _, r := range privateRanges {
		if bytes4(ip) >= bytes4(r.start) && bytes4(ip) <= bytes4(r.end) {
			return true
		}
	}

	// Check for IPv6 private addresses
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	return false
}

// Helper function to convert IPv4 to uint32 for range comparison
func bytes4(ip net.IP) uint32 {
	if len(ip) == 16 {
		// Convert IPv4-mapped IPv6 to IPv4
		ip = ip[12:16]
	}
	if len(ip) != 4 {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
