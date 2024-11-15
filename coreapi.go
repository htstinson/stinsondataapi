// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     Database // Interface for database operations
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Database interface for loose coupling
type Database interface {
	GetItem(id string) (*Item, error)
	CreateItem(item *Item) error
	// Add other methods as needed
}

func (app *App) Initialize() {
	app.Router = mux.NewRouter()
	app.setupRoutes()
	app.setupMiddleware()
}

func (app *App) setupRoutes() {
	// Group routes by version
	v1 := app.Router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/items", app.createItem).Methods("POST")
	v1.HandleFunc("/items/{id}", app.getItem).Methods("GET")
	v1.HandleFunc("/health", app.healthCheck).Methods("GET")
}

func (app *App) setupMiddleware() {
	app.Router.Use(loggingMiddleware)
	app.Router.Use(jsonContentTypeMiddleware)
}

// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

// Middleware to set JSON content type
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (app *App) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, Response{
		Status:  code,
		Message: message,
	})
}

func (app *App) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}

func (app *App) createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	item.CreatedAt = time.Now()

	if err := app.DB.CreateItem(&item); err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Error creating item")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, Response{
		Status: http.StatusCreated,
		Data:   item,
	})
}

func (app *App) getItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item, err := app.DB.GetItem(vars["id"])
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   item,
	})
}

func (app *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	app.respondWithJSON(w, http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Service is healthy",
	})
}

func main() {
	app := &App{}
	app.Initialize()

	srv := &http.Server{
		Handler:      app.Router,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server starting on port 8080")
	log.Fatal(srv.ListenAndServe())
}
