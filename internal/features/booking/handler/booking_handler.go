package bookinghandler

import (
	"encoding/json"
	"net/http"

	"github.com/codepnw/stdlib-ticket-system/internal/errs"
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

func (h *bookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	var req BookingCancelReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := utils.Validate(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.uc.CancelBooking(r.Context(), req.BookingID)
	if err != nil {
		switch err {
		case errs.ErrBookingNotFound:
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
		case errs.ErrCancelOtherBooking:
			helper.ErrorResponse(w, http.StatusConflict, err.Error())
		case errs.ErrBookingIsCancel:
			helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		case errs.ErrBookingIsPaid:
			helper.ErrorResponse(w, http.StatusBadRequest, err.Error())

		default:
			helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	helper.SuccessResponse(w, http.StatusOK, "booking cancelled", nil)
}
