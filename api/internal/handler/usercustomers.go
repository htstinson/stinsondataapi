package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectUserCustomerView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SelectUserCustomerView")

	ctx := r.Context()
	user_customer_views, err := h.db.SelectUserCustomerView(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select user_customer_view")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_customer_views)

}

func (h *Handler) GetUserCustomer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("h GetUserCustomer")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	user_customer, err := h.db.GetUserCustomer(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user_customer")
		return
	}
	if user_customer == nil {
		common.RespondError(w, http.StatusNotFound, "User_Customer not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_customer)
}

func (h *Handler) UpdateUserCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateUserCustomer")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var user_customer = model.User_Customer{}
	if err := json.NewDecoder(r.Body).Decode(&user_customer); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	current_user_customer, err := h.db.GetUserCustomer(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user_customer")
		return
	}

	if current_user_customer == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "User_Customer not found")
		return
	}

	current_user_customer.User_ID = user_customer.User_ID
	current_user_customer.Customer_Id = user_customer.Customer_Id

	err = h.db.UpdateUserCustomer(ctx, user_customer)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating customer")
		return
	}

	fmt.Println(5)

	common.RespondJSON(w, http.StatusOK, user_customer)
}

func (h *Handler) CreateUserCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateUserCustomer")

	var user_customer *model.User_Customer
	if err := json.NewDecoder(r.Body).Decode(&user_customer); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	_, err := h.db.LookupUserCustomer(ctx, user_customer.User_ID, user_customer.Customer_Id)
	if err != nil {
		if err.Error() == "not found" {
			// do nothing
		} else {
			fmt.Println(err.Error())
			fmt.Println("duplicate user customer")
			return
		}
	}

	fmt.Println("ok")

	new_user_customer, err := h.db.CreateUserCustomer(ctx, user_customer.User_ID, user_customer.Customer_Id)
	if err != nil {
		fmt.Println("Could not create user_customer")
		common.RespondError(w, http.StatusInternalServerError, "Failed to create user_customer")
		return
	}

	common.RespondJSON(w, http.StatusCreated, new_user_customer)
}

func (h *Handler) DeleteUserCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	user_customer, err := h.db.GetUserCustomer(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}
	if user_customer == nil {
		common.RespondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	err = h.db.DeleteUserCustomer(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting customer")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_customer)

}
