package handler

import (
	"net/http"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) SelectPermissions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	permissions, err := h.db.SelectPermissions(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select permissions")
		return
	}

	common.RespondJSON(w, http.StatusOK, permissions)

}
