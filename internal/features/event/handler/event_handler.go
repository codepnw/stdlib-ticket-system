package eventhandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/codepnw/stdlib-ticket-system/internal/errs"
	"github.com/codepnw/stdlib-ticket-system/internal/features/event"
	eventusecase "github.com/codepnw/stdlib-ticket-system/internal/features/event/usecase"
	"github.com/codepnw/stdlib-ticket-system/internal/helper"
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
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.uc.CreateEvent(r.Context(), req); err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(w, http.StatusCreated, "event created.", nil)
}

func (h *eventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	data, err := h.uc.GetAllEvents(r.Context())
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "", data)
}

func (h *eventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ParseInt64(r.PathValue("event_id"))
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.uc.GetEventByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrEventNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "", data)
}

func (h *eventHandler) GetSeatsByEventID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ParseInt64(r.PathValue("event_id"))
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.uc.GetSeatsByEventID(r.Context(), id)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "", data)
}
