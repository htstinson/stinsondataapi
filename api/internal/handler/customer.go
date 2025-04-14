package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Customer - Create, Update, Delete, Get, List

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer *model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	newcustomer, err := h.db.CreateCustomer(ctx, customer.Name)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	common.RespondJSON(w, http.StatusCreated, newcustomer)
}

func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateCustomer")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var customer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	fmt.Println("h UpdateCustomer", customer.CreatedAt, customer.Name)

	currentcustomer, err := h.db.GetCustomer(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}

	if currentcustomer == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	currentcustomer.Name = customer.Name

	err = h.db.UpdateCustomer(ctx, currentcustomer)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating customer")
		return
	}

	fmt.Println(5)

	common.RespondJSON(w, http.StatusOK, customer)
}

func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	customer, err := h.db.GetCustomer(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}
	if customer == nil {
		common.RespondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	err = h.db.DeleteCustomer(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting customer")
		return
	}

	common.RespondJSON(w, http.StatusOK, customer)

}

func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	customer, err := h.db.GetCustomer(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}
	if customer == nil {
		common.RespondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, customer)
}

func (h *Handler) SelectCustomers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	customers, err := h.db.SelectCustomers(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list customers")
		return
	}

	common.RespondJSON(w, http.StatusOK, customers)

}
