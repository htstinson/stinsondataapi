package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (h *Handler) SelectPermissions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	permissions, err := h.db.SelectPermissions(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select permissions")
		return
	}

	common.RespondJSON(w, http.StatusOK, permissions)

}
