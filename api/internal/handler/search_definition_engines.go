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

func (h *Handler) SelectSearchDefinitionEnginesView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelecttSearchDefinitionEnginesView")

	ctx := r.Context()

	vars := mux.Vars(r)
	subscriber_Id := vars["subscriber_id"]

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to select subcriber")
		return
	}

	results, err := h.db.SelectSearchDefinitionEnginesSubscriberView(ctx, *subscriber, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select search definition engines view")
		return
	}

	common.RespondJSON(w, http.StatusOK, results)

}

func (h *Handler) DeleteSearchDefinitionEngine(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteSearchDefinitionEngine")

	vars := mux.Vars(r)
	subscriber_id := vars["subscriber_id"]
	search_definition_engine_id := vars["search_definition_engine_id"]

	ctx := r.Context()

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	search_engine, err := h.db.GetSearchDefinitionEnginesView(ctx, *subscriber, search_definition_engine_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = h.db.DeleteSearchDefinitionEngine(ctx, subscriber, search_definition_engine_id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting Search Definition Engine")
		return
	}

	common.RespondJSON(w, http.StatusOK, search_engine)

}

func (h *Handler) CreateSearchDefinitionEngines(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateSearchDefinitionEngines")

	var row *model.SearchDefinitionEngines
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

	fmt.Println("id", row.Id)
	fmt.Println("row.SearchEngineId", row.SearchEngineId)
	fmt.Println("row.SearchDefinitionsId", row.SearchDefinitionsId)

	row, err = h.db.CreateSearchDefinitionEngine(ctx, *subcriber, *row)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create row")
		return
	}

	common.RespondJSON(w, http.StatusCreated, row)
}
