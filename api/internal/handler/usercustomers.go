package handler

import (
	"fmt"
	"net/http"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
)

func (h *Handler) SelectUserCustomerView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SelectUserCustomerView")

	ctx := r.Context()
	user_customer_views, err := h.db.SelectUserCustomerView(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select permissions")
		return
	}
	for k, v := range user_customer_views {
		fmt.Println(k, v.Id, v.User_ID, v.Customer_Id, v.User_Username, v.Customer_Name, v.Assignedd_At)
	}

	common.RespondJSON(w, http.StatusOK, user_customer_views)

}
