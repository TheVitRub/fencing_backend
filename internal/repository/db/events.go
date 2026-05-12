package db

import (
	"context"
	"fmt"
	"time"

	dbEntities "fencing-club/internal/domain/entities/db"
)

// ListEvents возвращает все события, отсортированные по дате убывания.
func (r *DBRepository) ListEvents(ctx context.Context) ([]dbEntities.Event, error) {
	var events []dbEntities.Event
	err := r.db.Select(ctx, nil, &events,
		`SELECT id, title, description, date, location, image_url, images, created_at
		 FROM events
		 ORDER BY date DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListEvents: %w", err)
	}
	return events, nil
}

// CreateEvent добавляет событие в БД и проставляет ID в переданную структуру.
func (r *DBRepository) CreateEvent(ctx context.Context, e *dbEntities.Event) error {
	if e.Images == "" {
		e.Images = "[]"
	}
	err := r.db.QueryRow(ctx, nil,
		[]any{&e.ID},
		`INSERT INTO events (title, description, date, location, image_url, images, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		e.Title, e.Description, e.Date, e.Location, e.ImageURL, e.Images, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("CreateEvent: %w", err)
	}
	return nil
}

// UpdateEvent обновляет поля события по ID.
func (r *DBRepository) UpdateEvent(ctx context.Context, e *dbEntities.Event) error {
	if e.Images == "" {
		e.Images = "[]"
	}
	_, err := r.db.Exec(ctx, nil,
		`UPDATE events
		 SET title=$1, description=$2, date=$3, location=$4, image_url=$5, images=$6
		 WHERE id=$7`,
		e.Title, e.Description, e.Date, e.Location, e.ImageURL, e.Images, e.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateEvent id=%d: %w", e.ID, err)
	}
	return nil
}

// DeleteEvent удаляет событие по ID.
func (r *DBRepository) DeleteEvent(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, nil, `DELETE FROM events WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("DeleteEvent id=%d: %w", id, err)
	}
	return nil
}
