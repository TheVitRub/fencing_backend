package db

import (
	"context"
	"fmt"
	"time"

	dbEntities "fencing-club/internal/domain/entities/db"
)

// GetPage возвращает страницу по slug.
// При отсутствии возвращает sql.ErrNoRows.
func (r *DBRepository) GetPage(ctx context.Context, slug string) (*dbEntities.Page, error) {
	var p dbEntities.Page
	err := r.db.Get(ctx, &p,
		`SELECT id, slug, title, content, updated_at FROM pages WHERE slug = $1`,
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("GetPage slug=%s: %w", slug, err)
	}
	return &p, nil
}

// UpsertPage создаёт страницу или обновляет существующую по slug (ON CONFLICT).
func (r *DBRepository) UpsertPage(ctx context.Context, slug, title, content string) error {
	_, err := r.db.Exec(ctx, nil,
		`INSERT INTO pages (slug, title, content, updated_at)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (slug) DO UPDATE
		   SET title      = EXCLUDED.title,
		       content    = EXCLUDED.content,
		       updated_at = EXCLUDED.updated_at`,
		slug, title, content, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("UpsertPage slug=%s: %w", slug, err)
	}
	return nil
}
