package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
	"github.com/htstinson/stinsondataapi/api/pkg/database"
	"github.com/htstinson/stinsondataapi/api/pkg/database/schema"
)

// Subscriber - Create, Update, Delete, Get, List

func (h *Handler) CreateSubscriber(w http.ResponseWriter, r *http.Request) {
	var subscriber *model.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subscriber); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	newsubscriber, err := h.db.CreateSubscriber(ctx, subscriber.Name)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create subscriber")
		return
	}

	db := h.db.(*database.Database).DB

	schema_name := fmt.Sprintf("%s_", strings.ToLower(newsubscriber.Name[:3]))

	schema_name += strings.ReplaceAll(newsubscriber.Id, "-", "_")

	schema := schema.Schema{
		DB:             db,
		FromSchemaName: "subscriber_template",
		ToSchemaName:   schema_name,
	}

	err = schema.CopySchema(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = h.db.CreateProfile(ctx, schema_name, newsubscriber.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	common.RespondJSON(w, http.StatusCreated, newsubscriber)
}

func (h *Handler) UpdateSubscriber(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateSubscriber")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var subscriber model.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subscriber); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	fmt.Println("h UpdateSubscriber", subscriber.CreatedAt, subscriber.Name)

	currentsubscriber, err := h.db.GetSubscriber(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	if currentsubscriber == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "subscriber not found")
		return
	}

	currentsubscriber.Name = subscriber.Name

	err = h.db.UpdateSubscriber(ctx, currentsubscriber)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating subscriber")
		return
	}

	fmt.Println(5)

	common.RespondJSON(w, http.StatusOK, subscriber)
}

func (h *Handler) DeleteSubscriber(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteSubscriber")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	subscriber, err := h.db.GetSubscriber(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}
	if subscriber == nil {
		common.RespondError(w, http.StatusNotFound, "subscriber not found")
		return
	}

	err = h.db.DeleteSubscriber(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting subscriber")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriber)

}

func (h *Handler) GetSubscriber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	subscriber, err := h.db.GetSubscriber(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}
	if subscriber == nil {
		common.RespondError(w, http.StatusNotFound, "subscriber not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriber)
}

func (h *Handler) SelectSubscribers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	subscribers, err := h.db.SelectSubscribers(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list subscribers")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscribers)

}

func (h *Handler) Create_Schema(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h Create_Schema")

	var subscriber *model.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subscriber); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	id := subscriber.Id
	if id == "" {
		common.RespondError(w, http.StatusBadRequest, "No Id")
		return
	}

	fmt.Println(id)

	ctx := r.Context()

	subscriber, err := h.db.GetSubscriber(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}
	if subscriber == nil {
		common.RespondError(w, http.StatusNotFound, "subscriber not found")
		return
	}

	db := h.db.(*database.Database).DB

	schema_name := fmt.Sprintf("%s_", strings.ToLower(subscriber.Name[:3]))

	schema_name += strings.ReplaceAll(subscriber.Id, "-", "_")

	schema := schema.Schema{
		DB:             db,
		FromSchemaName: "subscriber_template",
		ToSchemaName:   schema_name,
	}

	err = schema.CopySchema(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = h.db.CreateProfile(ctx, schema_name, subscriber.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

}
