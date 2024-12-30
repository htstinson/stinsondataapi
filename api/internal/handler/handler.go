package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"api/internal/auth"
	"api/internal/model"
	"api/internal/salesforce"
	"api/pkg/database"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db     database.Repository
	auth   auth.JWTAuth
	sfauth salesforce.SalesforceAuth
	logger *log.Logger
}

func NewHandler(db database.Repository, auth auth.JWTAuth, sfauth salesforce.SalesforceAuth, logger *log.Logger) *Handler {
	return &Handler{db: db, auth: auth, sfauth: sfauth, logger: logger}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Item - Create, Update, Delete, Get, List

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

// User - Create, Update, Delete, Get, List

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	password := "password"

	ctx := r.Context()
	user, err := h.db.CreateUser(ctx, user.Username, password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentuser, err := h.db.GetUser(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	if currentuser == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	currentuser.Username = user.Username
	err = h.db.UpdateUser(ctx, currentuser)
	if err != nil {
		respondError(w, http.StatusNotFound, "Error updating user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	user, err := h.db.GetUser(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	if user == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	err = h.db.DeleteUser(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Error deleting user")
		return
	}

	respondJSON(w, http.StatusOK, user)

}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	user, err := h.db.GetUser(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	if user == nil {
		respondError(w, http.StatusNotFound, "Item not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	items, err := h.db.ListUsers(ctx, 100, 0)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	respondJSON(w, http.StatusOK, items)

}

// Salesforce

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("GetAccount")
	vars := mux.Vars(r)
	id := vars["id"]

	account, err := salesforce.GetAccount(h.sfauth, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "account lookup error")
		return
	}

	if account.Id == "" {
		respondError(w, http.StatusNotFound, "account not found")
		return
	}

	h.logger.Println(account)

	respondJSON(w, http.StatusOK, account)
}

func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("UpdateAccount")
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Println("id: ", id)

	var newAccount model.NewAccount // this is for new or updated accounts

	if err := json.NewDecoder(r.Body).Decode(&newAccount); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentAccount, err := salesforce.GetAccount(h.sfauth, id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if currentAccount.Id == "" {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	endpoint := fmt.Sprintf("/services/data/v61.0/sobjects/Account/%s", currentAccount.Id)

	_, err = salesforce.SalesforcePatch(h.sfauth, endpoint, newAccount)
	if err != nil {
		respondError(w, http.StatusNotFound, "Error updating account")
		return
	}

	h.logger.Println(newAccount.Name, newAccount.AccountType, *newAccount.AccountSource)

	respondJSON(w, http.StatusOK, newAccount)
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	h.logger.Println("Create Account")

	var bodyBytes bytes.Buffer
	_, err := bodyBytes.ReadFrom(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	// Display the body
	log.Printf("Body: %s", bodyBytes.String())

	// Restore the body for further processing
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes.Bytes()))

	var account *model.NewAccount
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	fmt.Println(account)

	response, err := salesforce.SalesforcePost(h.sfauth, "/services/data/v62.0/sobjects/Account", account)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	h.logger.Println(string(response))

	respondJSON(w, http.StatusOK, "test complete")
}

func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {

	query := `SELECT 
	Id,
	Name,
	Industry,
	Description,
	Phone,
	Fax,
	Website,
	LastModifiedDate, 
	CreatedDate,
	LastActivityDate,	
	LastViewedDate,
	IsDeleted,
	MasterRecordId,
	Type,
	ParentId,
	BillingStreet,
	BillingCity,
	BillingState,
	BillingPostalCode,
	BillingCountry,
	AnnualRevenue,
	NumberOfEmployees,
	OwnerId,
	CreatedById,
	LastModifiedById,
	AccountSource
	FROM Account LIMIT 200`

	data, err := salesforce.SalesforceGet(h.sfauth, "/services/data/v59.0/query?q=", query, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	response := model.AccountQueryResponse{}

	err = json.Unmarshal(data, &response)
	if err != nil {
		h.logger.Println(err.Error())
	}

	respondJSON(w, http.StatusOK, response.Records)
}

// All - Internal
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

// All - UI
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
