package seat

type seatStatus string

const (
	StatusAvailable seatStatus = "AVAILABLE"
	StatusReserved  seatStatus = "RESERVED"
	StatusSold      seatStatus = "SOLD"
)

type Seat struct {
	ID         int64      `json:"id" db:"id"`
	EventID    int64      `json:"event_id" db:"event_id"`
	SeatNumber string     `json:"seat_number" db:"seat_number"`
	Price      float64    `json:"price" db:"price"`
	Status     seatStatus `json:"status" db:"status"`
	Version    int        `json:"version" db:"version"`
}
