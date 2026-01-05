package eventrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/stdlib-ticket-system/internal/errs"
	"github.com/codepnw/stdlib-ticket-system/internal/features/event"
)

type EventRepository interface {
	CreateEventTx(ctx context.Context, tx *sql.Tx, input event.Event) (int64, error)
	GetEventByID(ctx context.Context, eventID int64) (event.Event, error)
	GetAllEvents(ctx context.Context) ([]event.Event, error)
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

func (r *eventRepository) GetEventByID(ctx context.Context, eventID int64) (event.Event, error) {
	query := `
		SELECT id, name, event_date, is_active, created_at, updated_at
		FROM events WHERE id = $1 LIMIT 1
	`
	var e event.Event
	
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(
		&e.ID,
		&e.Name,
		&e.EventDate,
		&e.IsActive,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return event.Event{}, errs.ErrEventNotFound
		}
		return event.Event{}, err
	}
	return e, nil
}

func (r *eventRepository) GetAllEvents(ctx context.Context) ([]event.Event, error) {
	query := `
		SELECT id, name, event_date, is_active, created_at, updated_at 
		FROM events ORDER BY id DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var events []event.Event
	for rows.Next() {
		var e event.Event
		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.EventDate,
			&e.IsActive,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}