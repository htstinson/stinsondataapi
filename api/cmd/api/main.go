package main

import (
	"fmt"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/auth"
	"github.com/htstinson/stinsondataapi/api/internal/handler"
	"github.com/htstinson/stinsondataapi/api/internal/middleware"
	"github.com/htstinson/stinsondataapi/api/internal/model"
	"github.com/htstinson/stinsondataapi/api/pkg/database"

	//searcher "github.com/htstinson/business_searcher"

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

	//	fmt.Printf("[%v] [main] Initializing SalesForce.com connection.\n", time.Now().Format(time.RFC3339))

	//	sf, err := salesforce.New()
	//	if err != nil {
	//		fmt.Printf("[%v] [main] SalesForce error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
	//		return
	//	}

	fmt.Printf("[%v] [main] Initializing RDS database.\n", time.Now().Format(time.RFC3339))
	var RDSLogin = &model.RDSLogin{}
	rdsLogin, err := common.GetSecretString("RDS/apidb", "us-west-2")
	if err != nil {
		fmt.Printf("[%v] [main] RDS error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
		return
	}
	json.Unmarshal(rdsLogin, RDSLogin)

	// Initialize database
	config := database.Config{
		Host:        RDSLogin.Host,
		Port:        RDSLogin.Port,
		User:        RDSLogin.Username,
		Password:    RDSLogin.Password,
		Search_Path: "common",
		DBName:      "apidb",
		SSLMode:     "require",
	}

	db, err := database.New(config)
	if err != nil {
		fmt.Printf("[%v] [main] Failed to connect to RDS database: %s.\n", time.Now().Format(time.RFC3339), err.Error())
		return
	}
	defer db.Close()
	fmt.Printf("[%v] [main] Connected to RDS database.\n", time.Now().Format(time.RFC3339))

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

	//api.Use(middleware.IpLoggingMiddleware)

	// Public routes
	api.HandleFunc("/health", h.HealthCheck).Methods("GET")
	api.HandleFunc("/register", h.Register).Methods("POST")
	api.HandleFunc("/login", h.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/", h.HealthCheck).Methods("GET")

	// Protected routes
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(middleware.IpLoggingMiddleware)
	protected.Use(middleware.CORS) // First: Set CORS headers
	protected.Use(jwtAuth.Middleware)

	//test
	protected.HandleFunc("/test", h.Test).Methods("POST", "OPTIONS")

	// Blocked
	protected.HandleFunc("/blocked/update", h.AddBlockedFromRDSToWAF).Methods("GET", "OPTIONS")
	protected.HandleFunc("/blocked/parse", h.AddBlockedFromLogs).Methods("GET", "OPTIONS")
	protected.HandleFunc("/blocked/{id}", h.UpdateBlocked).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/blocked/{id}", h.GetBlocked).Methods("GET", "OPTIONS")
	protected.HandleFunc("/blocked/{id}", h.DeleteBlocked).Methods("DELETE")
	protected.HandleFunc("/blocked", h.CreateBlocked).Methods("POST", "OPTIONS")
	protected.HandleFunc("/blocked", h.SelectBlocked).Methods("GET", "OPTIONS")

	// Item
	protected.HandleFunc("/items", h.CreateItem).Methods("POST", "OPTIONS")
	protected.HandleFunc("/items/{id}", h.UpdateItem).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/items/{id}", h.GetItem).Methods("GET")
	protected.HandleFunc("/items", h.ListItems).Methods("GET", "OPTIONS")
	protected.HandleFunc("/items/{id}", h.DeleteItem).Methods("DELETE")

	// Salesforce Account
	//protected.HandleFunc("/accounts", sf.Handler.CreateAccount).Methods("POST", "OPTIONS")
	//protected.HandleFunc("/accounts/{id}", sf.Handler.UpdateAccount).Methods("PATCH", "OPTIONS")
	//protected.HandleFunc("/accounts", sf.Handler.ListAccounts).Methods("GET", "OPTIONS")

	// Salesforce Contact
	//protected.HandleFunc("/contacts", sf.Handler.ListContacts).Methods("GET", "OPTIONS")
	//protected.HandleFunc("/contacts/{accountid}", sf.Handler.ListContacts).Methods("GET", "OPTIONS")
	//protected.HandleFunc("/contact/{contactid}", sf.Handler.GetContactById).Methods("GET", "OPTIONS")

	// User
	protected.HandleFunc("/users", h.CreateUser).Methods("POST", "OPTIONS")
	protected.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
	protected.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	protected.HandleFunc("/users", h.SelectUsers).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/roles", h.SelectUserRoles).Methods("GET", "OPTIONS")
	protected.HandleFunc("/profile", h.GetUser).Methods("GET", "OPTIONS")

	// User_Subscriber
	protected.HandleFunc("/usersubscriberview/user/{id}", h.SelectUserSubscriberViewByUserId).Methods("GET", "OPTIONS")
	protected.HandleFunc("/usersubscriberview", h.SelectUserSubscriberView).Methods("GET", "OPTIONS")
	protected.HandleFunc("/usersubscriber/{id}", h.UpdateUserSubscriber).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/usersubscriber", h.CreateUserSubscriber).Methods("POST", "OPTIONS")
	protected.HandleFunc("/usersubscriber/{id}", h.DeleteUserSubscriber).Methods("DELETE")

	// UserSubscriberRole
	protected.HandleFunc("/usersubscriberroleview", h.SelectUserSubscriberRoleView).Methods("POST", "OPTIONS")
	protected.HandleFunc("/usersubscriberrole", h.CreateUserSubscriberRole).Methods("POST", "OPTIONS")
	protected.HandleFunc("/usersubscriberrole/{id}", h.UpdateUserSubscriberRole).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/usersubscriberrole/{id}", h.DeleteUserSubscriberRole).Methods("DELETE")

	// Subscriber - Customer
	protected.HandleFunc("/subscriber/customers", h.SelectSubscriberCustomers).Methods("POST", "OPTIONS")
	protected.HandleFunc("/subscriber/customer", h.CreateCustomer).Methods("POST", "OPTIONS")
	protected.HandleFunc("/subscriber/customerd", h.DeleteCustomer).Methods("POST", "OPTIONS")
	protected.HandleFunc("/subscriber/customer", h.UpdateCustomer).Methods("PUT", "OPTIONS")

	// Subscriber - Customer - Contacts
	protected.HandleFunc("/subscriber/customer/contacts", h.SelectContacts).Methods("POST", "OPTIONS")
	protected.HandleFunc("/subscriber/customer/contact", h.CreateContact).Methods("POST", "OPTIONS")
	protected.HandleFunc("/subscriber/customer/contactd", h.DeleteContact).Methods("POST", "OPTIONS")
	protected.HandleFunc("/subscriber/customer/contact", h.UpdateContact).Methods("PUT", "OPTIONS")

	// Subsriber - Item View
	protected.HandleFunc("/subscriber/items/{id}", h.SelectSubscriberItemView).Methods("GET", "OPTIONS")
	protected.HandleFunc("/subscriber/item/{id}", h.DeleteSubscriberItem).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/subscriber/item", h.CreateSubscriberItem).Methods("POST", "OPTIONS")

	// Subscribers
	protected.HandleFunc("/subscribers", h.CreateSubscriber).Methods("POST", "OPTIONS")
	protected.HandleFunc("/subscribers/{id}", h.UpdateSubscriber).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/subscribers", h.DeleteSubscriber).Methods("DELETE")
	protected.HandleFunc("/subscibers/{id}", h.GetSubscriber).Methods("GET")
	protected.HandleFunc("/subscribers/g", h.GetSubscriberP).Methods("POST", "OPTIONS")
	protected.HandleFunc("/subscribers", h.SelectSubscribers).Methods("GET", "OPTIONS")

	// Role
	protected.HandleFunc("/roles", h.CreateRole).Methods("POST", "OPTIONS")
	protected.HandleFunc("/roles/{id}", h.UpdateRole).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/roles/{id}", h.DeleteRole).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/roles/{id}", h.GetRole).Methods("GET")
	protected.HandleFunc("/roles", h.SelectRoles).Methods("GET", "OPTIONS")

	// Permission
	protected.HandleFunc("/permissions", h.CreatePermission).Methods("POST", "OPTIONS")
	protected.HandleFunc("/permissions/{id}", h.UpdatePermission).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/permissions/{id}", h.DeletePermission).Methods("DELETE")
	protected.HandleFunc("/permissions", h.SelectPermissions).Methods("GET", "OPTIONS")

	// User Permission

	// Role Permission
	protected.HandleFunc("/rolepermissionsview", h.SelectRolePermissionsView).Methods("GET", "OPTIONS")

	// Add middleware
	api.Use(middleware.Logger(&log.Logger{}))
	api.Use(middleware.RequestID)
	api.Use(middleware.SecurityHeaders)
	api.Use(middleware.CORS)

	//static assets
	distPath := "/home/ec2-user/go/src/stinsondata-tools-reactapp/dist"
	fmt.Printf("[%v] [main] Serving files from: %s.\n", time.Now().Format(time.RFC3339), distPath)

	// Handle all static assets including the index.js file
	router.PathPrefix("/assets/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%v] [main] Asset request for: %s.\n", time.Now().Format(time.RFC3339), r.URL.Path)

		// Remove the leading /assets/ to get the file path
		filePath := filepath.Join(distPath, r.URL.Path)
		fmt.Printf("[%v] [main] Looking for file at: %s.\n", time.Now().Format(time.RFC3339), filePath)

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
		fmt.Printf("[%v] [main] Server starting...\n", time.Now().Format(time.RFC3339))

		err := srv.ListenAndServeTLS("../../certs/certificate.crt", "../../certs/private.key")
		if err == http.ErrServerClosed {
			fmt.Printf("[%v] [main] Failed to start server (tls): %v.\n", time.Now().Format(time.RFC3339), err.Error())
		} else {
			fmt.Printf("[%v] [main] Error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
		}

	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	fmt.Printf("[%v] [main] Server stopping...\n", time.Now().Format(time.RFC3339))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("[%v] [main] Server forced to shutdown: %s.\n", time.Now().Format(time.RFC3339), err.Error())
		return
	}

	fmt.Printf("[%v] [main] Server stopped.\n", time.Now().Format(time.RFC3339))
}

/*
func test() {
	configFile := flag.String("config", "", "Path to a single YAML configuration file (optional)")
	configDir := flag.String("config-dir", "./search_definitions", "Path to directory containing YAML configuration files")
	outputDir := flag.String("output-dir", "./search_results", "Path to directory for output files")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found, using system environment variables")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("ERROR: GOOGLE_API_KEY environment variable is not set")
	}

	// Ensure output directory exists
	if err := searcher.EnsureDirectory(*outputDir); err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	// Determine which config files to process
	var configFiles []string
	var err error

	if *configFile != "" {
		// Single file mode (backward compatibility)
		// If the config file doesn't have a path separator, check in config-dir
		configPath := *configFile
		if !strings.Contains(configPath, string(filepath.Separator)) && !filepath.IsAbs(configPath) {
			// Just a filename, look in config-dir
			configPath = filepath.Join(*configDir, configPath)
			fmt.Printf("Looking for config file in: %s\n", configPath)
		}
		configFiles = []string{configPath}
		fmt.Printf("Processing single config file: %s\n", configPath)
	} else {
		// Directory mode (new default)
		fmt.Printf("Looking for config files in: %s\n", *configDir)
		configFiles, err = searcher.LoadConfigsFromDirectory(*configDir)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
		fmt.Printf("Found %d config file(s)\n", len(configFiles))
		for i, cf := range configFiles {
			fmt.Printf("  [%d] %s\n", i+1, cf)
		}
	}

	// Process each config file
	totalSearches := 0
	totalSuccessful := 0
	totalFailed := 0

	for _, cfgFile := range configFiles {
		fmt.Printf("\n--- Processing: %s ---\n", cfgFile)

		config, err := searcher.LoadConfig(cfgFile)
		if err != nil {
			log.Printf("ERROR: Failed to load configuration from '%s': %v", cfgFile, err)
			continue
		}

		client, err := searcher.NewSearchClient(apiKey, config)
		if err != nil {
			log.Printf("ERROR: Failed to create search client for '%s': %v", cfgFile, err)
			continue
		}

		// Generate output filename
		timestamp := time.Now().Format("20060102_150405")
		configName := strings.TrimSuffix(filepath.Base(cfgFile), filepath.Ext(cfgFile))
		outputFilename := filepath.Join(*outputDir, fmt.Sprintf("search_results_%s_%s.json", configName, timestamp))

		fmt.Printf("Executing searches...\n")

		// Build output structure
		output := searcher.OutputResult{
			Timestamp:     time.Now().Format(time.RFC3339),
			ConfigFile:    cfgFile,
			Configuration: client.BuildConfigurationOutput(),
			Searches:      client.ExecuteAllSearches(),
		}

		// Write JSON to file
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			log.Printf("ERROR: Failed to create output file '%s': %v", outputFilename, err)
			continue
		}

		encoder := json.NewEncoder(outputFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			outputFile.Close()
			log.Printf("ERROR: Failed to write JSON to '%s': %v", outputFilename, err)
			continue
		}
		outputFile.Close()

		// Calculate statistics
		successCount := 0
		for _, search := range output.Searches {
			if search.Error == "" {
				successCount++
			}
		}

		fmt.Printf("✓ Search results saved to: %s\n", outputFilename)
		fmt.Printf("  Total searches: %d\n", len(output.Searches))
		fmt.Printf("  Successful: %d\n", successCount)
		if successCount < len(output.Searches) {
			fmt.Printf("  Failed: %d\n", len(output.Searches)-successCount)
		}

		totalSearches += len(output.Searches)
		totalSuccessful += successCount
		totalFailed += (len(output.Searches) - successCount)
	}

	// Print summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Config files processed: %d\n", len(configFiles))
	fmt.Printf("Total searches executed: %d\n", totalSearches)
	fmt.Printf("Total successful: %d\n", totalSuccessful)
	if totalFailed > 0 {
		fmt.Printf("Total failed: %d\n", totalFailed)
	}
	fmt.Printf("Results directory: %s\n", *outputDir)
}
*/
