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

	fmt.Println(subscriber_Id)

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_Id)

	search_engines, err := h.db.SelectSearchEngines(ctx, *subscriber, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select user_subscriber_view")
		return
	}

	common.RespondJSON(w, http.StatusOK, search_engines)

}
