package db

import (
	"context"
	"fmt"
	"time"

	dbEntities "fencing-club/internal/domain/entities/db"
)

// GetFounder возвращает данные об основателе школы.
// В таблице всегда одна запись с id=1.
// Если запись ещё не создана — возвращает sql.ErrNoRows.
func (r *DBRepository) GetFounder(ctx context.Context) (*dbEntities.Founder, error) {
	var f dbEntities.Founder
	err := r.db.Get(ctx, &f,
		`SELECT id, name, bio, photo_url, updated_at FROM founder LIMIT 1`,
	)
	if err != nil {
		return nil, fmt.Errorf("GetFounder: %w", err)
	}
	return &f, nil
}

// UpsertFounder создаёт или обновляет единственную запись об основателе.
// id жёстко равен 1 — в таблице может быть только один основатель.
func (r *DBRepository) UpsertFounder(ctx context.Context, f *dbEntities.Founder) error {
	_, err := r.db.Exec(ctx, nil,
		`INSERT INTO founder (id, name, bio, photo_url, updated_at)
		 VALUES (1, $1, $2, $3, $4)
		 ON CONFLICT (id) DO UPDATE
		   SET name       = EXCLUDED.name,
		       bio        = EXCLUDED.bio,
		       photo_url  = EXCLUDED.photo_url,
		       updated_at = EXCLUDED.updated_at`,
		f.Name, f.Bio, f.PhotoURL, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("UpsertFounder: %w", err)
	}
	return nil
}
