package bookinghandler

import (
	"encoding/json"
	"net/http"

	bookingusecase "github.com/codepnw/stdlib-ticket-system/internal/features/booking/usecase"
	"github.com/codepnw/stdlib-ticket-system/internal/helper"
	"github.com/codepnw/stdlib-ticket-system/pkg/utils"
)

type bookingHandler struct {
	uc bookingusecase.BookingUsecase
}

func NewBookingHandler(uc bookingusecase.BookingUsecase) *bookingHandler {
	return &bookingHandler{uc: uc}
}

func (h *bookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req BookingCreateReq
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := utils.Validate(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	if err := h.uc.CreateBooking(r.Context(), req.EventID, req.SeatIDs); err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(w, http.StatusOK, "event booked", nil)
}

func (h *bookingHandler) GetBookingHistory(w http.ResponseWriter, r *http.Request) {
	data, err := h.uc.GetBookingHistory(r.Context())
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(w, http.StatusOK, "", data)
}