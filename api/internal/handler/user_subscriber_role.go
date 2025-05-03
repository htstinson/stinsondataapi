package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectUserSubscriberRolesView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectUserSubscriberRolesView")

	ctx := r.Context()
	user_customer_roles_views, err := h.db.SelectUserSubscriberRolesView(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select user_customer_roles_view")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_customer_roles_views)

}

func (h *Handler) CreateUserSubscriberRole(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateUserSubscriberRole")

	var user_subscriber_role *model.User_Subscriber_Role

	if err := json.NewDecoder(r.Body).Decode(&user_subscriber_role); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	fmt.Println("user_subscriber_id", user_subscriber_role.User_Subscriber_ID)
	fmt.Println("role_id", user_subscriber_role.Role_Id)
	fmt.Println(&user_subscriber_role)

	_, err := h.db.LookupUserSubscriberRole(ctx, user_subscriber_role.User_Subscriber_ID, user_subscriber_role.Role_Id)
	if err != nil {
		if err.Error() == "not found" {
			// do nothing
		} else {
			fmt.Println(err.Error())
			fmt.Println("duplicate user subscriber role")
			return
		}
	}

	fmt.Println("ok")

	new_user_subscriber_role, err := h.db.CreateUserSubscriberRole(ctx, user_subscriber_role.User_Subscriber_ID, user_subscriber_role.Role_Id)
	if err != nil {
		fmt.Println("Could not create user_subscriber_role")
		common.RespondError(w, http.StatusInternalServerError, "Failed to create user_subscribe_role")
		return
	}

	common.RespondJSON(w, http.StatusCreated, new_user_subscriber_role)
}
