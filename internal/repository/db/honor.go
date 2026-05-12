package db

import (
	"context"
	"fmt"
	"time"

	dbEntities "fencing-club/internal/domain/entities/db"
)

// ListHonorMembers возвращает участников доски почёта, отсортированных по sort_order.
func (r *DBRepository) ListHonorMembers(ctx context.Context) ([]dbEntities.HonorMember, error) {
	var members []dbEntities.HonorMember
	err := r.db.Select(ctx, nil, &members,
		`SELECT id, name, title, description, photo_url, sort_order, created_at
		 FROM honor_members
		 ORDER BY sort_order ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListHonorMembers: %w", err)
	}
	return members, nil
}

// CreateHonorMember добавляет участника на доску почёта и проставляет ID.
func (r *DBRepository) CreateHonorMember(ctx context.Context, m *dbEntities.HonorMember) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&m.ID},
		`INSERT INTO honor_members (name, title, description, photo_url, sort_order, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		m.Name, m.Title, m.Description, m.PhotoURL, m.SortOrder, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("CreateHonorMember: %w", err)
	}
	return nil
}

// UpdateHonorMember обновляет поля участника доски почёта по ID.
func (r *DBRepository) UpdateHonorMember(ctx context.Context, m *dbEntities.HonorMember) error {
	_, err := r.db.Exec(ctx, nil,
		`UPDATE honor_members
		 SET name=$1, title=$2, description=$3, photo_url=$4, sort_order=$5
		 WHERE id=$6`,
		m.Name, m.Title, m.Description, m.PhotoURL, m.SortOrder, m.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateHonorMember id=%d: %w", m.ID, err)
	}
	return nil
}

// DeleteHonorMember удаляет участника доски почёта по ID.
func (r *DBRepository) DeleteHonorMember(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, nil, `DELETE FROM honor_members WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("DeleteHonorMember id=%d: %w", id, err)
	}
	return nil
}
