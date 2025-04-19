package handler

import (
	"fmt"
	"net/http"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) SelectUserCustomerRolesView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SelectUserCustomerRolesView")

	ctx := r.Context()
	user_customer_roles_views, err := h.db.SelectUserCustomerRolesView(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select user_customer_roles_view")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_customer_roles_views)

}
