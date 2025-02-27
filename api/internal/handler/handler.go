package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/auth"
	"github.com/htstinson/stinsondataapi/api/internal/model"
	"github.com/htstinson/stinsondataapi/api/pkg/database"

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
	fmt.Printf("[%v] HealthCheck\n", time.Now().Format(time.RFC3339))
	common.RespondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Item - Create, Update, Delete, Get, List

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := h.db.CreateItem(ctx, &item); err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	common.RespondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentitem, err := h.db.GetItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if currentitem == nil {
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}

	currentitem.Name = item.Name
	err = h.db.UpdateItem(ctx, currentitem)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error updating item")
		return
	}

	common.RespondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	item, err := h.db.GetItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if item == nil {
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}
	err = h.db.DeleteItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting item")
		return
	}

	common.RespondJSON(w, http.StatusOK, item)

}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	item, err := h.db.GetItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if item == nil {
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := h.db.ListItems(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	common.RespondJSON(w, http.StatusOK, items)
}

// Admin

func (h *Handler) ListBlocked(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ListBlocked")
	ctx := r.Context()
	items, err := h.db.ListBlocked(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	fmt.Printf("Items = %v\n", len(items))

	common.RespondJSON(w, http.StatusOK, items)
}

// User - Create, Update, Delete, Get, List

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	password := "password"

	ctx := r.Context()
	user, err := h.db.CreateUser(ctx, user.Username, password)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	common.RespondJSON(w, http.StatusCreated, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentuser, err := h.db.GetUser(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	if currentuser == nil {
		common.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	currentuser.Username = user.Username
	err = h.db.UpdateUser(ctx, currentuser)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error updating user")
		return
	}

	common.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	user, err := h.db.GetUser(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	if user == nil {
		common.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	err = h.db.DeleteUser(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting user")
		return
	}

	common.RespondJSON(w, http.StatusOK, user)

}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	user, err := h.db.GetUser(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	if user == nil {
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	items, err := h.db.ListUsers(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	common.RespondJSON(w, http.StatusOK, items)

}

// All - UI
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate input
	if req.Username == "" || req.Password == "" {
		common.RespondError(w, http.StatusBadRequest, "Username and password required")
		return
	}

	// Check if user exists
	existingUser, err := h.db.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Error checking username")
		return
	}
	if existingUser != nil {
		common.RespondError(w, http.StatusConflict, "Username already exists")
		return
	}

	// Create user
	user, err := h.db.CreateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	common.RespondJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%v] Login\n", time.Now().Format(time.RFC3339))
	var req model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("[%v] Error: %s\n", time.Now().Format(time.RFC3339), err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Get user
	user, err := h.db.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		fmt.Printf("[%v] Error: %s\n", time.Now().Format(time.RFC3339), err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Error finding user")
		return
	}
	if user == nil {
		common.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		fmt.Printf("[%v] Error: %s\n", time.Now().Format(time.RFC3339), err.Error())
		common.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token
	token, err := h.auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		fmt.Printf("[%v] Error: %s\n", time.Now().Format(time.RFC3339), err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	common.RespondJSON(w, http.StatusOK, model.LoginResponse{
		Token:     token,
		ExpiresIn: int64(h.auth.Config.TokenDuration.Seconds()),
	})
}
