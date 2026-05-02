package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectCustomers")

	sort := ""
	order := ""

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

	customers, _, err := h.db.SelectCustomers(ctx, *subcriber, 100, 0, sort, order)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select customers")
		return
	}

	common.RespondJSON(w, http.StatusOK, customers)

}

func (h *Handler) SelectSubscriberCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectSubscriberCustomers")

	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	// Sensible defaults if missing/zero
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var subcriber *model.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subcriber); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	subcriber, err := h.db.GetSubscriber(ctx, subcriber.Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
	}

	customers, total, err := h.db.SelectCustomers(ctx, *subcriber, limit, offset, sort, order)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select customers")
		return
	}

	//common.RespondJSON(w, http.StatusOK, customers)

	common.RespondJSON2(w, http.StatusOK, map[string]any{
		"data":  customers,
		"total": total,
	})

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

	subscriber, err := h.db.GetSubscriber(ctx, customer.Subscriber_Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
	}

	profile, err := h.db.GetProfile(ctx, subscriber)

	customer.Id = uuid.New().String()
	customer.Subscriber_Id = subscriber.Id
	customer.Schema_Name = subscriber.Schema_Name

	customer, err = h.db.CreateCustomer(ctx, customer, profile)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	common.RespondJSON(w, http.StatusCreated, customer)
}

func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteCustomer")

	vars := mux.Vars(r)
	id := vars["customer_id"]
	subscriber_id := vars["subscriber_id"]
	ctx := r.Context()

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
	}

	var customer = model.Customer{
		Id:            id,
		Subscriber_Id: subscriber_id,
		Schema_Name:   subscriber.Schema_Name,
	}

	current, err := h.db.GetCustomer(ctx, customer)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}
	if current == nil {
		fmt.Println("Customer not found")
		common.RespondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	contacts, err := h.db.SelectContacts(ctx, *current, 100, 0)

	if err != nil {
		fmt.Println(err.Error())
		if err != sql.ErrNoRows {
			common.RespondError(w, http.StatusOK, "Error locating contacts")
			return
		}
	}

	if len(contacts) > 0 {
		common.RespondError(w, http.StatusConflict, "Cannot delete customer: customer has associated contacts")
		return
	}

	err = h.db.DeleteCustomer(ctx, &customer)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting Customer")
		return
	}

	common.RespondJSON(w, http.StatusOK, customer)

}

func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateCustomer")
	ctx := r.Context()

	var customer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	fmt.Println(customer.Name)

	current, err := h.db.GetCustomer(ctx, customer)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get customer")
		return
	}
	if current == nil {
		common.RespondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	err = h.db.UpdateCustomer(ctx, &customer)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error updating customer")
		return
	}

	common.RespondJSON(w, http.StatusOK, customer)
}
