package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) ListSearchEngines(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h ListtSearchEngines")

	ctx := r.Context()

	vars := mux.Vars(r)
	subscriber_Id := vars["subscriber_id"]

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to select subcriber")
		return
	}

	search_engines, err := h.db.SelectSearchEngines(ctx, *subscriber, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select search engines")
		return
	}

	common.RespondJSON(w, http.StatusOK, search_engines)

}
