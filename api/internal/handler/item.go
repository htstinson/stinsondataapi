package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Item - Create, Update, Delete, Get, List

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := h.db.CreateItem(ctx, &item); err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	common.RespondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	currentitem, err := h.db.GetItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if currentitem == nil {
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}

	currentitem.Name = item.Name
	err = h.db.UpdateItem(ctx, currentitem)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error updating item")
		return
	}

	common.RespondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()

	item, err := h.db.GetItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if item == nil {
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}
	err = h.db.DeleteItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting item")
		return
	}

	common.RespondJSON(w, http.StatusOK, item)

}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	item, err := h.db.GetItem(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}
	if item == nil {
		common.RespondError(w, http.StatusNotFound, "Item not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := h.db.SelectItems(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	common.RespondJSON(w, http.StatusOK, items)
}
