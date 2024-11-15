package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"myapi/internal/auth"
	"myapi/internal/model"
	"myapi/pkg/database"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db   database.Repository
	auth auth.JWTAuth
}

func NewHandler(db database.Repository, auth auth.JWTAuth) *Handler {
	return &Handler{db: db, auth: auth}
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

	respondJSON(w, http.StatusOK, model.LoginResponse{
		Token:     token,
		ExpiresIn: int64(h.auth.Config.TokenDuration.Seconds()),
	})
}
