package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
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
