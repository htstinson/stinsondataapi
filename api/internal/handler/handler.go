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
	fmt.Printf("[%v] [Login]\n", time.Now().Format(time.RFC3339))
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

	fmt.Println("Login user", user.IP_address)

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
	roles, err := h.db.SelectRolesByUser(r.Context(), user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(user.ID, roles)

	token, err := h.auth.GenerateToken(*user, roles)
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
