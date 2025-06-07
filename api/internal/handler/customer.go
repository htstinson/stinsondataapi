package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectCustomers")

	var user *model.CurrentUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	subcriber, err := h.db.GetSubscriber(ctx, user.Subscribed[0].Subscriber_ID)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
	}

	customers, err := h.db.SelectCustomers(ctx, *subcriber, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list customers")
		return
	}

	common.RespondJSON(w, http.StatusOK, customers)

}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateCustomer")

	var customer *model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	ctx := r.Context()

	subcriber, err := h.db.GetSubscriber(ctx, customer.Subscriber_ID)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
	}

	customer.Id = uuid.New().String()
	customer.Subscriber_ID = subcriber.Id
	customer.Schema_Name = subcriber.Schema_Name

	customer, err = h.db.CreateCustomer(ctx, customer)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	common.RespondJSON(w, http.StatusCreated, customer)
}

func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteCustomer")

	var customer *model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	current, err := h.db.GetCustomer(ctx, *&customer)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}
	if current == nil {
		common.RespondError(w, http.StatusNotFound, "Contact not found")
		return
	}
	err = h.db.DeleteCustomer(ctx, customer)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting item")
		return
	}

	common.RespondJSON(w, http.StatusOK, customer)

}
