package handler

import (
	"fmt"
	"net/http"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) SelectRolePermissions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	role_permissions, err := h.db.SelectRolePermissions(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select permissions")
		return
	}
	for k, v := range role_permissions {
		fmt.Println(k, v.Role_Id, v.V_Role_Name, v.Permission_Id, v.V_Permission_Name, v.Object_Id, v.V_Object_Name, v.V_Object_Type)
	}

	common.RespondJSON(w, http.StatusOK, role_permissions)

}
