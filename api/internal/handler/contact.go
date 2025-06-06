package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

	contacts, err := h.db.SelectContacts(ctx, *customer, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list contacts")
		return
	}

	common.RespondJSON(w, http.StatusOK, contacts)

}

func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateContact")

	var contact *model.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	contact.Id = uuid.New().String()
	fmt.Println(contact.Schema_Name_)
	fmt.Println(contact.ParentId)

	contact, err := h.db.CreateContact(ctx, contact)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	common.RespondJSON(w, http.StatusCreated, contact)
}
