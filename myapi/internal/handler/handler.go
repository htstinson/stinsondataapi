package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"myapi/internal/auth"
	"myapi/internal/model"
	"myapi/internal/salesforce"
	"myapi/pkg/database"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db     database.Repository
	auth   auth.JWTAuth
	logger *log.Logger
}

func NewHandler(db database.Repository, auth auth.JWTAuth, logger *log.Logger) *Handler {
	return &Handler{db: db, auth: auth, logger: logger}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := h.db.CreateItem(ctx, &item); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentitem, err := h.db.GetItem(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if currentitem == nil {
		respondError(w, http.StatusNotFound, "Item not found")
		return
	}

	currentitem.Name = item.Name
	err = h.db.UpdateItem(ctx, currentitem)
	if err != nil {
		respondError(w, http.StatusNotFound, "Error updating item")
		return
	}

	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	item, err := h.db.GetItem(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Item not found")
		return
	}

	fmt.Println(item)

	err = h.db.DeleteItem(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Error deleting item")
		return
	}

	respondJSON(w, http.StatusOK, item)

}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Println(id)

	ctx := r.Context()
	item, err := h.db.GetItem(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if item == nil {
		respondError(w, http.StatusNotFound, "Item not found")
		return
	}

	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := h.db.ListItems(ctx, 100, 0)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) Account(w http.ResponseWriter, r *http.Request) {

	authResponse, err := SalesForceLogin()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	auth := salesforce.SalesforceAuth{
		AccessToken: authResponse.AccessToken,
		InstanceURL: "https://stinsondata.my.salesforce.com",
	}

	account := model.Account{
		Name:        "A Test Account",
		Description: "Created via API",
		Phone:       "1234567890",
	}

	// TODO: This should be async
	// Question: Does each user have a unique session?

	response, err := salesforce.SalesforcePost(auth, "/services/data/v62.0/sobjects/Account", account)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	h.logger.Println(string(response))

	respondJSON(w, http.StatusOK, "test complete")
}

func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {

	authResponse, err := SalesForceLogin()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	auth := salesforce.SalesforceAuth{
		AccessToken: authResponse.AccessToken,
		InstanceURL: "https://stinsondata.my.salesforce.com",
	}

	query := "SELECT Id, Name, Industry FROM Account LIMIT 1000"

	data, err := salesforce.SalesforceGet(auth, "/services/data/v59.0/query?q=", query, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	response := model.AccountQueryResponse{}

	err = json.Unmarshal(data, &response)
	if err != nil {
		// Handle error
	}

	for k, v := range response.Records {
		fmt.Println(k, v.Name, v.Phone, v.Id)
	}

	respondJSON(w, http.StatusOK, "get complete")
}

func SalesForceLogin() (*salesforce.SalesforceAuthResponse, error) {
	// Your Salesforce credentials
	var (
		clientID     = os.Getenv("SF_STINSONDATA_CLIENT_ID")
		clientSecret = os.Getenv("SF_STINSONDATA_CLIENT_SECRET")
		username     = os.Getenv("SF_STINSONDATA_USERNAME")
		password     = os.Getenv("SF_STINSONDATA_PASSWORD")
		loginURL     = "https://login.salesforce.com"
	)

	auth, err := salesforce.GetSalesforceToken(clientID, clientSecret, username, password, loginURL)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return nil, err
	}

	return auth, err

}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate input
	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password required")
		return
	}

	// Check if user exists
	existingUser, err := h.db.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error checking username")
		return
	}
	if existingUser != nil {
		respondError(w, http.StatusConflict, "Username already exists")
		return
	}

	// Create user
	user, err := h.db.CreateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Get user
	user, err := h.db.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error finding user")
		return
	}
	if user == nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token
	token, err := h.auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	respondJSON(w, http.StatusOK, model.LoginResponse{
		Token:     token,
		ExpiresIn: int64(h.auth.Config.TokenDuration.Seconds()),
	})
}
