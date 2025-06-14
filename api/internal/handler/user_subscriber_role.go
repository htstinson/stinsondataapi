package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectUserSubscriberRolesView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelectUserSubscriberRolesView")

	ctx := r.Context()
	user_customer_roles_views, err := h.db.SelectUserSubscriberRoleView(ctx, 100, 0)
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

func (h *Handler) UpdateUserSubscriberRole(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h UpdateUserSubscriberRole")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var user_subscriber_role = model.User_Subscriber_Role{}

	if err := json.NewDecoder(r.Body).Decode(&user_subscriber_role); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	fmt.Println("id", user_subscriber_role.Id)
	fmt.Println("role id", user_subscriber_role.Role_Id)
	fmt.Println("user_subscriber_id", user_subscriber_role.User_Subscriber_ID)

	current_user_subscriber_role, err := h.db.GetUserSubscriberRole(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user_subscriber_role")
		return
	}

	if current_user_subscriber_role == nil {
		fmt.Println(3)
		common.RespondError(w, http.StatusNotFound, "User_Subscriber not found")
		return
	}

	err = h.db.UpdateUserSubscriberRole(ctx, user_subscriber_role)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating user subscriber role")
		return
	}

	fmt.Println(5)

	common.RespondJSON(w, http.StatusOK, user_subscriber_role)
}

func (h *Handler) DeleteUserSubscriberRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	user_subscriber_role, err := h.db.GetUserSubscriberRole(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get user_subscriber_role")
		return
	}
	if user_subscriber_role == nil {
		common.RespondError(w, http.StatusNotFound, "User_Subscriber_Role not found")
		return
	}

	err = h.db.DeleteUserSubscriberRole(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting user_subscriber_role")
		return
	}

	common.RespondJSON(w, http.StatusOK, user_subscriber_role)

}
