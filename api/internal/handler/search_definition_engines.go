package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
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
