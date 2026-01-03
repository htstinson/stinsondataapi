package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (h *Handler) SelectSearchEngines(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h SelecttSearchEngines")

	ctx := r.Context()

	vars := mux.Vars(r)
	subscriber_Id := vars["subscriber_id"]

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_Id)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to select subcriber")
		return
	}

	search_engines, err := h.db.SelectSearchEngines(ctx, *subscriber, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to select search engines")
		return
	}

	common.RespondJSON(w, http.StatusOK, search_engines)

}

func (h *Handler) CreateSearchEngine(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateSearchEngine")

	var search_engine *model.SearchEngine
	if err := json.NewDecoder(r.Body).Decode(&search_engine); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	ctx := r.Context()

	search_engine.Id = uuid.New().String()

	subcriber, err := h.db.GetSubscriber(ctx, search_engine.SubscriberId)
	if err != nil {
		fmt.Println(err.Error())
		common.RespondError(w, http.StatusInternalServerError, "Failed to get subscriber")
	}

	search_engine, err = h.db.CreateSearchEngine(ctx, *search_engine, *subcriber)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to create search engine")
		return
	}

	common.RespondJSON(w, http.StatusCreated, search_engine)
}

func (h *Handler) DeleteSearchEngine(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h DeleteSearchEngine")

	var search_engine *model.SearchEngine
	if err := json.NewDecoder(r.Body).Decode(&search_engine); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	subscriber, err := h.db.GetSubscriber(ctx, search_engine.SubscriberId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = h.db.DeleteSearchEngine(ctx, *search_engine, subscriber)
	if err != nil {
		common.RespondError(w, http.StatusNotFound, "Error deleting Search Engine")
		return
	}

	common.RespondJSON(w, http.StatusOK, search_engine)

}
