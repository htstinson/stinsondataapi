package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) SelectSearchResults(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectSearchResults")

	ctx := r.Context()

	vars := mux.Vars(r)
	subscriber_Id := vars["subscriber_id"]
	searchDefinitionEngineId := vars["search_definition_engine_id"]

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to select subcriber")
		return
	}

	results, err := h.db.SelectSearchResultView(ctx, *subscriber, searchDefinitionEngineId)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select search definitions")
		return
	}

	common.RespondJSON(w, http.StatusOK, results)

}
