package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreatePermission")
	var permission *model.Permission
	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	newpermission, err := h.db.CreatePermission(ctx, permission.Name, permission.Description)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create permission")
		return
	}

	common.RespondJSON(w, http.StatusCreated, newpermission)
}

func (h *Handler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdatePermission")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var permission model.Permission
	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentpermission, err := h.db.GetPermission(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get permission")
		return
	}

	if currentpermission == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "Permission not found")
		return
	}

	currentpermission.Name = permission.Name
	currentpermission.Description = permission.Description

	err = h.db.UpdatePermission(ctx, currentpermission)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating permission")
		return
	}

	fmt.Println(5)

	common.RespondJSON(w, http.StatusOK, permission)
}

func (h *Handler) SelectPermissions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	permissions, err := h.db.SelectPermissions(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select permissions")
		return
	}

	common.RespondJSON(w, http.StatusOK, permissions)

}
