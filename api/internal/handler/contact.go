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

func (h *Handler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteContact")
	var contact *model.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	contact, err := h.db.GetContact(ctx, contact)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get contact")
		return
	}
	if contact == nil {
		common.RespondError(w, http.StatusNotFound, "Contact not found")
		return
	}
	err = h.db.DeleteContact(ctx, contact)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting item")
		return
	}

	common.RespondJSON(w, http.StatusOK, contact)

}

func (h *Handler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpateContact")
	ctx := r.Context()

	var contact model.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	current, err := h.db.GetContact(ctx, &contact)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get contact")
		return
	}
	if current == nil {
		common.RespondError(w, http.StatusNotFound, "Contact not found")
		return
	}

	fmt.Println(contact)

	err = h.db.UpdateContact(ctx, &contact)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error updating contact")
		return
	}

	common.RespondJSON(w, http.StatusOK, contact)
}
