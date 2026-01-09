package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectSearchDefinitions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectSearchDefinitions")

	ctx := r.Context()

	vars := mux.Vars(r)
	subscriber_Id := vars["subscriber_id"]

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to select subcriber")
		return
	}

	results, err := h.db.SelectSearchDefinitions(ctx, *subscriber, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select search definitions")
		return
	}

	common.RespondJSON(w, http.StatusOK, results)

}

func (h *Handler) DeleteSearchDefinition(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteSearchDefinition")

	vars := mux.Vars(r)
	subscriber_id := vars["subscriber_id"]
	search_definition_id := vars["search_definition_id"]

	ctx := r.Context()

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	search_engine, err := h.db.GetSearchDefinition(ctx, *subscriber, search_definition_id, 1, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = h.db.DeleteSearchDefinition(ctx, subscriber, search_definition_id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting Search Engine")
		return
	}

	common.RespondJSON(w, http.StatusOK, search_engine)

}

func (h *Handler) CreateSearchDefinition(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateSearchDefinition")

	var row *model.SearchDefinition
	if err := json.NewDecoder(r.Body).Decode(&row); err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	ctx := r.Context()

	row.Id = uuid.New().String()

	subcriber, err := h.db.GetSubscriber(ctx, row.SubscriberId)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
	}

	row, err = h.db.CreateSearchDefinition(ctx, *subcriber, *row)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create row")
		return
	}

	common.RespondJSON(w, http.StatusCreated, row)
}
