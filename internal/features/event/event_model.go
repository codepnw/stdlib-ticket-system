package event

import (
	"time"
)

type Event struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	EventDate time.Time `json:"event_date" db:"event_date"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SeatZoneReq struct {
	ZoneName    string  `json:"zone_name"`
	SeatsPerRow int     `json:"seats_per_row"`
	Price       float64 `json:"price"`
}

type CreateEventReq struct {
	Name      string        `json:"name"`
	EventDate time.Time     `json:"event_date"`
	IsActive  bool          `json:"is_active"`
	Zones     []SeatZoneReq `json:"zones"`
}
