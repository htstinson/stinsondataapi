package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectContacts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectContacts")

	var customer *model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	contacts, err := h.db.SelectContacts(ctx, customer.Schema_Name, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list contacts")
		return
	}

	common.RespondJSON(w, http.StatusOK, contacts)

}
