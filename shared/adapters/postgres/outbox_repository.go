package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yourusername/twitter-services-app/shared/domain"
)
type OutboxRepository struct {
	db *pgx.Conn
}

func NewOutboxRepository(db *pgx.Conn) *OutboxRepository {
	return &OutboxRepository{db: db}
}
func (r *OutboxRepository) Insert(
	ctx context.Context,
	event domain.OutboxEvent,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO outbox_events
		(id, event_type, aggregate_id, payload, correlation_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		event.ID,
		event.EventType,
		event.AggregateID,
		event.Payload,
		event.CorrelationID,
		event.CreatedAt,
	)

	return err
}

func (r *OutboxRepository) FetchPending(
	ctx context.Context,
	limit int,
) ([]domain.OutboxEvent, error) {

	rows, err := r.db.Query(ctx, `
		SELECT id, event_type, aggregate_id, payload, correlation_id, created_at
		FROM outbox_events
		WHERE sent_at IS NULL
		ORDER BY created_at
		LIMIT $1
	`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.OutboxEvent

	for rows.Next() {
		var e domain.OutboxEvent
		err := rows.Scan(
			&e.ID,
			&e.EventType,
			&e.AggregateID,
			&e.Payload,
			&e.CorrelationID,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}
func (r *OutboxRepository) MarkSent(
	ctx context.Context,
	eventID string,
) error {
	now := time.Now().UTC()

	_, err := r.db.Exec(ctx, `
		UPDATE outbox_events
		SET sent_at = $1
		WHERE id = $2
	`, now, eventID)

	return err
}
