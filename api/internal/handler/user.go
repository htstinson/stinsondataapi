package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/auth"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

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
		common.RespondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	common.RespondJSON(w, http.StatusCreated, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateUser")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	fmt.Println("h UpdateUser", user.CreatedAt, user.Username, user.IP_address)

	currentuser, err := h.db.GetUser(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if currentuser == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	currentuser.Username = user.Username
	currentuser.IP_address = user.IP_address
	err = h.db.UpdateUser(ctx, currentuser)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating user")
		return
	}

	fmt.Println(5)

	common.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) GetUserByUserName(w http.ResponseWriter, r *http.Request) {

	fmt.Println("h GetUserByUserName")

	var user *model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	user, err := h.db.GetUserByUsername(ctx, user.Username)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if user == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "Item not found")
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

	fmt.Println("h GetUser (no parms)")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	if id == "" {
		claims, ok := ctx.Value("user").(*auth.Claims)
		if !ok {
			fmt.Println("Type assertion failed: anyValue is not of type auth.Claims")
			return
		}
		id = claims.UserID
	}

	user, err := h.db.GetUser(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if user == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) SelectUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	items, err := h.db.SelectUsers(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	common.RespondJSON(w, http.StatusOK, items)

}

func (h *Handler) SelectUserRoles(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	items, err := h.db.SelectUserRoles(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	common.RespondJSON(w, http.StatusOK, items)

}
