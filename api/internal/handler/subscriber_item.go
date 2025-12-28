package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) SelectSubscriberItemView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SelectSubscriberItemView")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	subscriber_item_views, err := h.db.SelectSubscriberItemView(ctx, id, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select user_subscriber_view")
		return
	}

	common.RespondJSON(w, http.StatusOK, subscriber_item_views)

}
