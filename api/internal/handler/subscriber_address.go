package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectSubscriberAddresses(w http.ResponseWriter, r *http.Request) {
	// TODO
	fmt.Println("h Select Subscriber Addresses")

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
		return
	}

	addresses, total, err := h.db.SelectSubscriberAddresses(ctx, *subcriber, limit, offset, sort, order)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to select addresses")
		return
	}

	common.RespondJSON2(w, http.StatusOK, map[string]any{
		"data":  addresses,
		"total": total,
	})
}

func (h *Handler) GetSubscriberAddress(w http.ResponseWriter, r *http.Request) {
	// TODO
	fmt.Println("h Get Subscriber Address")

	var address *model.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	subcriber, err := h.db.GetSubscriber(ctx, address.SubscriberId)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	address, err = h.db.GetSubscriberAddress(ctx, subcriber.Schema_Name, address.Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to select addresses")
		return
	}

	common.RespondJSON(w, http.StatusOK, address)
}

func (h *Handler) UpdateSubscriberAddress(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateSubscriberAddress")
	ctx := r.Context()

	var address model.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	subscriber, err := h.db.GetSubscriber(ctx, address.SubscriberId)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	err = h.db.UpdateSubscriberAddress(ctx, subscriber, address)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to update subscriber address")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriber)
}

func (h *Handler) CreateSubscriberAddress(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateSubscriberAddress")
	ctx := r.Context()

	var address model.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	subscriber, err := h.db.GetSubscriber(ctx, address.SubscriberId)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	err = h.db.CreateSubscriberAddress(ctx, subscriber, address)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to update subscriber address")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriber)
}

func (h *Handler) DeleteSubscriberAddress(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h Delete Subscriber Address")

	vars := mux.Vars(r)
	id := vars["subscriber_id"]

	fmt.Println(id)

	ctx := r.Context()

	subscriberitem, err := h.db.GetSubscriberItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber_item")
		return
	}
	if subscriberitem == nil {
		common.RespondError(w, http.StatusNotFound, "Subscriber_item not found")
		return
	}

	err = h.db.DeleteSubscriberItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting subscriber_item")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriberitem)

}
