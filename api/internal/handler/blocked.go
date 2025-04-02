package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	waf "github.com/htstinson/stinsondataapi/api/aws/mywaf"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// blocked

func (h *Handler) ListBlocked(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	items, err := h.db.ListBlocked(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	waf.Block("Blocked", "", "", "us-west-2")

	common.RespondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateBlocked(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	ctx := r.Context()

	var blocked model.Blocked
	if err := json.NewDecoder(r.Body).Decode(&blocked); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	current, err := h.db.GetBlocked(ctx, id)
	if err != nil {
		fmt.Println(2, err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get blokced")
		return
	}
	if current == nil {
		fmt.Println(3, id, "not found")
		common.RespondError(w, http.StatusNotFound, "Blocked not found")
		return
	}

	current.IP = blocked.IP
	current.Notes = blocked.Notes
	err = h.db.UpdateBlocked(ctx, current)
	if err != nil {
		fmt.Println(4, err.Error())
		common.RespondError(w, http.StatusNotFound, "Error updating blocked")
		return
	}

	common.RespondJSON(w, http.StatusOK, blocked)
}

func (h *Handler) GetBlocked(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	user, err := h.db.GetBlocked(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get blocked")
		return
	}
	if user == nil {
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) CreateBlocked(w http.ResponseWriter, r *http.Request) {
	var blocked *model.Blocked
	if err := json.NewDecoder(r.Body).Decode(&blocked); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	newblocked, err := h.db.CreateBlocked(ctx, *blocked)
	if err != nil {
		fmt.Println("h create blocked ", err.Error())
		if err.Error() == "duplicate" {
			fmt.Println("h create blocked duplicate address")
			common.RespondJSON(w, 409, nil)
		}
		common.RespondError(w, http.StatusInternalServerError, "Failed to create blocked")
		return
	}

	waf.Block("Blocked", *&blocked.IP, "", "us-west-2")

	common.RespondJSON(w, http.StatusCreated, newblocked)
}

func (h *Handler) DeleteBlocked(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	blocked, err := h.db.GetBlocked(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get blocked")
		return
	}
	if blocked == nil {
		common.RespondError(w, http.StatusNotFound, "Blocked not found")
		return
	}

	err = h.db.DeleteBlocked(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting blocked")
		return
	}

	waf.Block("Blocked", "", blocked.IP, "us-west-2")

	common.RespondJSON(w, http.StatusOK, blocked)

}
