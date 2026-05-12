package db

import (
	"context"
	"fmt"
	"time"

	dbEntities "fencing-club/internal/domain/entities/db"
)

// ListPlans возвращает все учебные планы, отсортированные по дате создания убывания.
func (r *DBRepository) ListPlans(ctx context.Context) ([]dbEntities.Plan, error) {
	var plans []dbEntities.Plan
	err := r.db.Select(ctx, nil, &plans,
		`SELECT id, title, description, period, items, created_at
		 FROM plans
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListPlans: %w", err)
	}
	return plans, nil
}

// CreatePlan добавляет учебный план в БД и проставляет ID.
func (r *DBRepository) CreatePlan(ctx context.Context, p *dbEntities.Plan) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&p.ID},
		`INSERT INTO plans (title, description, period, items, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		p.Title, p.Description, p.Period, p.Items, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("CreatePlan: %w", err)
	}
	return nil
}

// UpdatePlan обновляет поля учебного плана по ID.
func (r *DBRepository) UpdatePlan(ctx context.Context, p *dbEntities.Plan) error {
	_, err := r.db.Exec(ctx, nil,
		`UPDATE plans
		 SET title=$1, description=$2, period=$3, items=$4
		 WHERE id=$5`,
		p.Title, p.Description, p.Period, p.Items, p.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdatePlan id=%d: %w", p.ID, err)
	}
	return nil
}

// DeletePlan удаляет учебный план по ID.
func (r *DBRepository) DeletePlan(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, nil, `DELETE FROM plans WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("DeletePlan id=%d: %w", id, err)
	}
	return nil
}
