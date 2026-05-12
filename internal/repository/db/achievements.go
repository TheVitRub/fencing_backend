package db

import (
	"context"
	"fmt"
	"time"

	dbEntities "fencing-club/internal/domain/entities/db"
)

// ListAchievements возвращает достижения школы, отсортированные по году убывания.
func (r *DBRepository) ListAchievements(ctx context.Context) ([]dbEntities.Achievement, error) {
	var list []dbEntities.Achievement
	err := r.db.Select(ctx, nil, &list,
		`SELECT id, title, description, year, image_url, created_at
		 FROM achievements
		 ORDER BY year DESC, id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListAchievements: %w", err)
	}
	return list, nil
}

// CreateAchievement добавляет достижение в БД и проставляет ID.
func (r *DBRepository) CreateAchievement(ctx context.Context, a *dbEntities.Achievement) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&a.ID},
		`INSERT INTO achievements (title, description, year, image_url, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		a.Title, a.Description, a.Year, a.ImageURL, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("CreateAchievement: %w", err)
	}
	return nil
}

// UpdateAchievement обновляет поля достижения по ID.
func (r *DBRepository) UpdateAchievement(ctx context.Context, a *dbEntities.Achievement) error {
	_, err := r.db.Exec(ctx, nil,
		`UPDATE achievements
		 SET title=$1, description=$2, year=$3, image_url=$4
		 WHERE id=$5`,
		a.Title, a.Description, a.Year, a.ImageURL, a.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateAchievement id=%d: %w", a.ID, err)
	}
	return nil
}

// DeleteAchievement удаляет достижение по ID.
func (r *DBRepository) DeleteAchievement(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, nil, `DELETE FROM achievements WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("DeleteAchievement id=%d: %w", id, err)
	}
	return nil
}
