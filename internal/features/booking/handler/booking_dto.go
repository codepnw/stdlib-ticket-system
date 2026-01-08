package bookinghandler

type BookingCreateReq struct {
	EventID int64   `json:"event_id" validate:"required"`
	SeatIDs []int64 `json:"seat_ids" validate:"required"`
}
