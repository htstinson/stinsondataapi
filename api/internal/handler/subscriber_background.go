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

func (h *Handler) SelectSubscriberBackgrounds(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h Select Subscriber Backgrounds")

	order := h.ValidOrder(r)
	sort := h.ValidSort(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var subscriber *model.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subscriber); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	fmt.Println("subscriber.Id", subscriber.Id)

	subscriber, err := h.db.GetSubscriber(ctx, subscriber.Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	backgrounds, total, err := h.db.SelectSubscriberBackgrounds(ctx, *subscriber, limit, offset, sort, order)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to select backgrounds")
		return
	}

	common.RespondJSON2(w, http.StatusOK, map[string]any{
		"data":  backgrounds,
		"total": total,
	})
}

func (h *Handler) GetSubscriberBackground(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h Get Subscriber Background")

	var background *model.Background
	if err := json.NewDecoder(r.Body).Decode(&background); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	subscriber, err := h.db.GetSubscriber(ctx, background.SubscriberId)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	fmt.Println("schema", subscriber.Schema_Name)
	fmt.Println("background.Id", background.Id)
	background, err = h.db.GetSubscriberBackground(ctx, subscriber.Schema_Name, background.Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get background")
		return
	}

	common.RespondJSON(w, http.StatusOK, background)
}

func (h *Handler) UpdateSubscriberBackground(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateSubscriberBackground")
	ctx := r.Context()

	var background model.Background
	if err := json.NewDecoder(r.Body).Decode(&background); err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	subscriber, err := h.db.GetSubscriber(ctx, background.SubscriberId)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	err = h.db.UpdateSubscriberBackground(ctx, subscriber, background)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to update subscriber background")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriber)
}

func (h *Handler) CreateSubscriberBackground(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateSubscriberBackground")
	ctx := r.Context()

	var background model.Background
	if err := json.NewDecoder(r.Body).Decode(&background); err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	subscriber, err := h.db.GetSubscriber(ctx, background.SubscriberId)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	err = h.db.CreateSubscriberBackground(ctx, subscriber, background)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create subscriber background")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriber)
}

func (h *Handler) DeleteSubscriberBackground(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h Delete Subscriber Background")

	vars := mux.Vars(r)
	subscriber_id := vars["subscriber_id"]
	background_id := vars["background_id"]

	ctx := r.Context()

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}
	if subscriber == nil {
		common.RespondError(w, http.StatusNotFound, "Subscriber not found")
		return
	}

	err = h.db.DeleteSubscriberBackground(ctx, subscriber.Schema_Name, background_id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting subscriber background")
		return
	}

	common.RespondJSON(w, http.StatusOK, nil)
}
