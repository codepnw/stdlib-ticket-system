package eventhandler

import (
	"encoding/json"
	"net/http"

	"github.com/codepnw/stdlib-ticket-system/internal/features/event"
	eventusecase "github.com/codepnw/stdlib-ticket-system/internal/features/event/usecase"
	"github.com/codepnw/stdlib-ticket-system/pkg/utils"
)

type eventHandler struct {
	uc eventusecase.EventUsecase
}

func NewEventHandler(uc eventusecase.EventUsecase) *eventHandler {
	return &eventHandler{uc: uc}
}

func (h *eventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req event.CreateEventReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateEvent(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "event created."})
}
