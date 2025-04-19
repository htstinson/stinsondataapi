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
		common.RespondError(w, http.StatusInternalServerError, "Failed to select permissions")
		return
	}
	for k, v := range user_customer_views {
		fmt.Println(k, v.Id, v.User_ID, v.Customer_Id, v.User_Username, v.Customer_Name, v.Assignedd_At)
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
