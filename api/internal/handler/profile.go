package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) GetSubscriberProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h GetSubscriberProfile")

	var subcriber *model.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subcriber); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	subcriber, err := h.db.GetSubscriber(ctx, subcriber.Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	profile, err := h.db.GetProfile(ctx, subcriber)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	common.RespondJSON(w, http.StatusOK, profile)
}
