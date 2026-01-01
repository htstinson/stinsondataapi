package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectSubscriberItemView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SelectSubscriberItemView")

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println(id)

	ctx := r.Context()
	subscriber_item_views, err := h.db.SelectSubscriberItemView(ctx, id, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select user_subscriber_view")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriber_item_views)

}

func (h *Handler) CreateSubscriberItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateSubscriberItem")

	var subscriber_item *model.Subscriber_Item
	if err := json.NewDecoder(r.Body).Decode(&subscriber_item); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	_, err := h.db.LookupSubscriberItem(ctx, subscriber_item.Item_ID, subscriber_item.Subscriber_Id)
	if err != nil {
		if err.Error() == "not found" {
			fmt.Println("subscriber item not found")
			// do nothing
		} else {
			fmt.Println(err.Error())
			fmt.Println("duplicate subscriber item")
			return
		}
	}

	fmt.Println("ok")

	new_user_subscriber, err := h.db.CreateSubscriberItem(ctx, subscriber_item.Item_ID, subscriber_item.Subscriber_Id)
	if err != nil {
		fmt.Println("Could not create subscriber_item")
		common.RespondError(w, http.StatusInternalServerError, "Failed to create subscriber_item")
		return
	}

	common.RespondJSON(w, http.StatusCreated, new_user_subscriber)
}

func (h *Handler) DeleteSubscriberItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteSubscriberItem")

	vars := mux.Vars(r)
	id := vars["id"]

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

func (h *Handler) GetSubscriberItem(w http.ResponseWriter, r *http.Request) {

	fmt.Println("h GetSubscriberItem")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	subscriberitem, err := h.db.GetSubscriberItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriberitem")
		return
	}
	if subscriberitem == nil {
		common.RespondError(w, http.StatusNotFound, "Subscriber_Item not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriberitem)
}
