package eventrepo

import (
	"context"
	"database/sql"

	"github.com/codepnw/stdlib-ticket-system/internal/features/event"
)

type EventRepository interface {
	CreateEventTx(ctx context.Context, tx *sql.Tx, input event.Event) (int64, error)
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEventTx(ctx context.Context, tx *sql.Tx, input event.Event) (int64, error) {
	query := `
		INSERT INTO events (name, event_date, is_active)
		VALUES ($1, $2, $3) RETURNING id
	`
	var eventID int64
	err := tx.QueryRowContext(
		ctx,
		query,
		input.Name,
		input.EventDate,
		input.IsActive,
	).Scan(&eventID)
	
	if err != nil {
		return 0, err
	}
	return eventID, nil
}

