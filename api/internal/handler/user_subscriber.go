package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectUserSubscriberView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SelectUserSubscriberView")

	ctx := r.Context()
	user_subscriber_views, err := h.db.SelectUserSubscriberView(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select user_subscriber_view")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_subscriber_views)

}

func (h *Handler) GetUserSubscriber(w http.ResponseWriter, r *http.Request) {

	fmt.Println("h GetUserSubscriber")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	user_subscriber, err := h.db.GetUserSubscriber(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user_subscriber")
		return
	}
	if user_subscriber == nil {
		common.RespondError(w, http.StatusNotFound, "User_Subscriber not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_subscriber)
}

func (h *Handler) UpdateUserSubscriber(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateUserSubscriber")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var user_subscriber = model.User_Subscriber{}
	if err := json.NewDecoder(r.Body).Decode(&user_subscriber); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	current_user_subscriber, err := h.db.GetUserSubscriber(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user_subscriber")
		return
	}

	if current_user_subscriber == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "User_Subscriber not found")
		return
	}

	current_user_subscriber.User_ID = user_subscriber.User_ID
	current_user_subscriber.Subscriber_Id = user_subscriber.Subscriber_Id

	err = h.db.UpdateUserSubscriber(ctx, user_subscriber)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating subscriber")
		return
	}

	fmt.Println(5)

	common.RespondJSON(w, http.StatusOK, user_subscriber)
}

func (h *Handler) CreateUserSubscriber(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateUserSubscriber")

	var user_subscriber *model.User_Subscriber
	if err := json.NewDecoder(r.Body).Decode(&user_subscriber); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	_, err := h.db.LookupUserSubscriber(ctx, user_subscriber.User_ID, user_subscriber.Subscriber_Id)
	if err != nil {
		if err.Error() == "not found" {
			// do nothing
		} else {
			fmt.Println(err.Error())
			fmt.Println("duplicate user subscriber")
			return
		}
	}

	fmt.Println("ok")

	new_user_subscriber, err := h.db.CreateUserSubscriber(ctx, user_subscriber.User_ID, user_subscriber.Subscriber_Id)
	if err != nil {
		fmt.Println("Could not create user_subscriber")
		common.RespondError(w, http.StatusInternalServerError, "Failed to create user_subscriber")
		return
	}

	common.RespondJSON(w, http.StatusCreated, new_user_subscriber)
}

func (h *Handler) DeleteUserSubscriber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	user_subscriber, err := h.db.GetUserSubscriber(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
		return
	}
	if user_subscriber == nil {
		common.RespondError(w, http.StatusNotFound, "Subscriber not found")
		return
	}

	err = h.db.DeleteUserSubscriber(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting subscriber")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_subscriber)

}
