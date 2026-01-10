package booking

import "time"

type bookingStatus string

const (
	StatusPending   bookingStatus = "PENDING"
	StatusPaid      bookingStatus = "PAID"
	StatusCancelled bookingStatus = "CANCELLED"
	StatusFailed    bookingStatus = "FAILED"
)

type Booking struct {
	ID          string        `json:"id" db:"id"`
	UserID      int64         `json:"user_id" db:"user_id"`
	EventID     int64         `json:"event_id" db:"event_id"`
	TotalAmount float64       `json:"total_amount" db:"total_amount"`
	Status      bookingStatus `json:"status" db:"status"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	ExpiresAt   time.Time     `json:"expires_at" db:"expires_at"`
}

type BookingItem struct {
	BookingID string `db:"booking_id"`
	SeatID    int64  `db:"seat_id"`
}

type BookingHistoryResponse struct {
	ID          string    `json:"id" db:"booking_id"`
	EventName   string    `json:"event_name" db:"event_name"`
	EventDate   time.Time `json:"event_date" db:"event_date"`
	TotalAmount float64   `json:"total_amount" db:"total_amount"`
	Status      string    `json:"status" db:"status"`
	SeatNumbers string    `json:"seat_numbers" db:"seat_numbers"` // STRING_AGG()
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
