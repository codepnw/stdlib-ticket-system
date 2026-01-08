package seat

type SeatStatus string

const (
	StatusAvailable SeatStatus = "AVAILABLE"
	StatusReserved  SeatStatus = "RESERVED"
	StatusSold      SeatStatus = "SOLD"
)

type Seat struct {
	ID         int64      `json:"id" db:"id"`
	EventID    int64      `json:"event_id" db:"event_id"`
	SeatNumber string     `json:"seat_number" db:"seat_number"`
	Price      float64    `json:"price" db:"price"`
	Status     SeatStatus `json:"status" db:"status"`
	Version    int        `json:"version" db:"version"`
}
