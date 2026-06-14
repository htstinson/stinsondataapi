package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/htstinson/stinsondataapi/api/commonweb"
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

	fmt.Println("Login user.IP", user.IP_address)

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
	// TODO replace this with roles per user_subscription
	roles, err := h.db.SelectRolesByUser(r.Context(), user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var user_subscriber_view = model.User_Subscriber_View{
		User_ID: user.ID,
	}

	fmt.Println("handler login user.ID = ", user.ID)

	user_subscriber_role_view, err := h.db.SelectUserSubscriberRoleView(r.Context(), user_subscriber_view, 100, 0)
	if err != nil {
		fmt.Println(err.Error())
	}

	token, err := h.auth.GenerateToken(*user, roles, user_subscriber_role_view)
	if err != nil {
		fmt.Printf("[%v] Error: %s\n", time.Now().Format(time.RFC3339), err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	login_time := time.Now().Format(time.RFC3339)

	body := fmt.Sprintf(
		`Your account logged into Thousand Hills Digital at %v.

If this was not you, please contact support@stinsondata.com.

Thank you!`, login_time)

	region := "us-west-2"
	/* TODO - Create a profile setting to determine if this is sent.
	    Need to have a choice of accounts to send from
		Also need to create the GCP connection to gmail automatically
		Also need to profile the body and subject of the message
		Full functionality will need to be able to read a mailbox
		Need e-mail address verification
		Need unsubscribe functionality
		Need to log e-mails that are sent
	*/
	if false {
		commonweb.SendMail(user.Username, "Thousand Hills Digital - Login", body, region)
	}

	common.RespondJSON(w, http.StatusOK, model.LoginResponse{
		Token:     token,
		ExpiresIn: int64(h.auth.Config.TokenDuration.Seconds()),
	})
}

func (h *Handler) ValidSort(r *http.Request) string {
	fmt.Println("h ValidSort")

	allowed := map[string]bool{"name": true, "username": true, "ip_address": true, "created_at": true}
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "id"
	}

	fmt.Println("sort = ", sort)

	if allowed[sort] {
		return sort
	}
	return sort
}

func (h *Handler) ValidOrder(r *http.Request) string {
	fmt.Println("h ValidOrder")

	order := r.URL.Query().Get("order")
	if order == "" {
		order = "asc"
	}

	fmt.Println("order = ", order)

	if order == "desc" {
		return "DESC"
	}
	return "ASC"
}
