package bookingrepo

import (
	"context"
	"database/sql"

	"github.com/codepnw/stdlib-ticket-system/internal/features/booking"
	"github.com/lib/pq"
)

//go:generate mockgen -source=booking_repo.go -destination=booking_repo_mock.go -package=bookingrepo
type BookingRepository interface {
	GetHistory(ctx context.Context, userID int64) ([]booking.BookingHistoryResponse, error)
	
	// Transaction
	CreateBookingTx(ctx context.Context, tx *sql.Tx, input booking.Booking) (string, error)
	CreateBookingItemsTx(ctx context.Context, tx *sql.Tx, bookingID string, seatIDs []int64) error
}

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateBookingTx(ctx context.Context, tx *sql.Tx, input booking.Booking) (string, error) {
	query := `
		INSERT INTO bookings (user_id, event_id, total_amount, status)
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	var id string

	err := tx.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.EventID,
		input.TotalAmount,
		input.Status,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *bookingRepository) CreateBookingItemsTx(ctx context.Context, tx *sql.Tx, bookingID string, seatIDs []int64) error {
	query := `
		INSERT INTO booking_items (booking_id, seat_id)
		SELECT $1, UNNEST($2::BIGINT[])
	`
	_, err := tx.ExecContext(ctx, query, bookingID, pq.Array(seatIDs))
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) GetHistory(ctx context.Context, userID int64) ([]booking.BookingHistoryResponse, error) {
	query := `
		SELECT
			b.id AS booking_id,
			e.name AS event_name,
			e.event_date,
			b.total_amount,
			b.status,
			b.created_at,
			STRING_AGG(s.seat_number, ', ') AS seat_numbers
		FROM bookings b
		JOIN events e ON b.event_id = e.id
		JOIN booking_items bi ON bi.booking_id = b.id 
		JOIN seats s ON bi.seat_id = s.id
		WHERE b.user_id = $1
		GROUP BY b.id, e.id
		ORDER BY b.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var history []booking.BookingHistoryResponse
	
	for rows.Next() {
		var h booking.BookingHistoryResponse
		if err := rows.Scan(
			&h.ID,
			&h.EventName,
			&h.EventDate,
			&h.TotalAmount,
			&h.Status,
			&h.CreatedAt,
			&h.SeatNumbers,
		); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return history, nil
}
