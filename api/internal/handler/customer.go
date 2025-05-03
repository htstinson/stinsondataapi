package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) SelectCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectCustomers")
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["schema_id"]

	customers, err := h.db.SelectCustomers(ctx, id, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list customers")
		return
	}

	common.RespondJSON(w, http.StatusOK, customers)

}
