package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) SelectCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h ListCustomers")
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	customers, err := h.db.LookupUserSubscribersByUserId(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list customers")
		return
	}

	common.RespondJSON(w, http.StatusOK, customers)

}
