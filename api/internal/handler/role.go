package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectRoles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectRoles")
	ctx := r.Context()
	customers, err := h.db.SelectRoles(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list customers")
		return
	}

	common.RespondJSON(w, http.StatusOK, customers)

}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateRole")

	var role *model.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	role, err := h.db.CreateRole(ctx, role.Name)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create role")
		return
	}

	common.RespondJSON(w, http.StatusCreated, role)
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateRole")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var role model.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentrole, err := h.db.GetRole(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get role")
		return
	}

	if currentrole == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "Role not found")
		return
	}

	currentrole.Name = role.Name

	err = h.db.UpdateRole(ctx, currentrole)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating role")
		return
	}

	fmt.Println(5)

	common.RespondJSON(w, http.StatusOK, role)
}

func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteRole")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	role, err := h.db.GetRole(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get role")
		return
	}
	if role == nil {
		common.RespondError(w, http.StatusNotFound, "Role not found")
		return
	}

	err = h.db.DeleteRole(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting fole")
		return
	}

	common.RespondJSON(w, http.StatusOK, role)

}

func (h *Handler) GetRole(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h GetRole")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	role, err := h.db.GetRole(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get fole")
		return
	}
	if role == nil {
		common.RespondError(w, http.StatusNotFound, "Role not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, role)
}
