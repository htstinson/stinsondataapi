package handler

import (
	"fmt"
	"net/http"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
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
