package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/htstinson/stinsondataapi/api/aws/mywaf"

	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
	"github.com/htstinson/stinsondataapi/api/internal/parser"
)

// blocked

func (h *Handler) ListBlocked(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	items, err := h.db.ListBlocked(ctx, 100, 0)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	mywaf.Block("Blocked", "", "", "us-west-2")

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
	fmt.Println("h GetBlocked (w,r)")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	blocked, err := h.db.GetBlocked(ctx, id)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, "Failed to get blocked")
		return
	}
	if blocked == nil {
		common.RespondError(w, http.StatusNotFound, "Item blocked not found")
		return
	}

	common.RespondJSON(w, http.StatusOK, blocked)
}

func (h *Handler) CreateBlocked(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h CreateBlocked(w,r)")
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

	mywaf.Block("Blocked", blocked.IP, "", "us-west-2")

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

	fmt.Println(blocked.IP)
	mywaf.Block("Blocked", "", blocked.IP, "us-west-2")

	common.RespondJSON(w, http.StatusOK, blocked)

}

func (h *Handler) AddBlockedFromLogs(w http.ResponseWriter, r *http.Request) {

	// Create blocked IP addresses from entries in the log.
	fmt.Printf("[%v] [main] Parse the log.\n", time.Now().Format(time.RFC3339))
	addresses, err := parser.ExtractUniqueIPsFromHandshakeErrors("/var/log/webserver.log")
	if err != nil {
		fmt.Printf("[%v] [main] error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
	} else {
		ctx := r.Context()
		fmt.Printf("[%v] [main] Blocked IP addresses.\n", time.Now().Format(time.RFC3339))
		for k, v := range addresses {
			blocked := &model.Blocked{
				Notes:     "TLS handshake error",
				CreatedAt: time.Now(),
			}
			ip := fmt.Sprintf("%s/32", v)
			blocked.IP = ip
			_, err := h.db.CreateBlocked(ctx, *blocked)
			if err == nil {
				fmt.Printf("[%v] [main] %v %s Created blocked IP.\n", time.Now().Format(time.RFC3339), k, ip)
			} else {
				fmt.Printf("[%v] [main] %s error: %s.\n", time.Now().Format(time.RFC3339), ip, err.Error())
			}

			err = mywaf.Block("Blocked", ip, "", "us-west-2")
			if err != nil {
				fmt.Printf("[%v] [main] %v %s Error adding IP to WAF IP Set.\n", time.Now().Format(time.RFC3339), k, ip)
			} else {
				fmt.Printf("[%v] [main] %v %s Added IP to WAF IP Set.\n", time.Now().Format(time.RFC3339), k, ip)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

}

func (h *Handler) AddBlockedFromRDSToWAF(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h AddBlockedFromRDSToWAF(w,r)")
	go func() {

		// Create blocked IP addresses from entries in RDS.
		fmt.Printf("[%v] [main] Add Blocked From RDS TO WAF.\n", time.Now().Format(time.RFC3339))
		ctx := context.Background()

		addresses, err := h.db.ListBlocked(ctx, 1000, 0)

		if err != nil {
			fmt.Printf("[%v] [main] error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
		} else {

			fmt.Printf("[%v] [main] Blocked %v IP addresses.\n", time.Now().Format(time.RFC3339), len(addresses))
			for k, v := range addresses {
				err = mywaf.Block("Blocked", v.IP, "", "us-west-2")
				if err != nil {
					fmt.Printf("[%v] [main] %v %s Error adding IP to WAF IP Set. %s\n", time.Now().Format(time.RFC3339), k, v.IP, err.Error())
				} else {
					fmt.Printf("[%v] [main] %v %s Added IP to WAF IP Set.\n", time.Now().Format(time.RFC3339), k, v.IP)
				}

				time.Sleep(500 * time.Millisecond)
			}
		}

	}()

	common.RespondJSON(w, http.StatusOK, "udating WAF from RDS")

}
