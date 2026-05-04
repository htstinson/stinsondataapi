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

func (h *Handler) UpdateSubscriberProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateSubscriberProfile")

	ctx := r.Context()

	var profile model.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	fmt.Println("54", profile.Id, *profile.Legal_Name, *profile.LinkedIn, profile.Subscriber_Id)

	var subscriber, err = h.db.GetSubscriber(ctx, profile.Subscriber_Id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}

	p, err := h.db.GetProfile(ctx, subscriber)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	err = h.db.UpdateProfile(ctx, p)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error updating subscriber")
		return
	}

	fmt.Println(p.Id, p.Legal_Name, p.LinkedIn)

	common.RespondJSON(w, http.StatusOK, subscriber)
}
