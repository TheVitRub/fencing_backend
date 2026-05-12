package application

import (
	"context"

	"fencing-club/internal/domain/apperrors"
	dbEntities "fencing-club/internal/domain/entities/db"
	"fencing-club/internal/domain/dto"
)

// ListPlans возвращает все учебные планы, отсортированные по дате создания убывания.
func (s *Service) ListPlans(ctx context.Context) ([]dto.PlanResponse, error) {
	plans, err := s.repo.ListPlans(ctx)
	if err != nil {
		return nil, wrapRepoError("ListPlans", err)
	}
	resp := make([]dto.PlanResponse, 0, len(plans))
	for _, p := range plans {
		resp = append(resp, dto.PlanResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Period:      p.Period,
			Items:       p.Items,
			CreatedAt:   p.CreatedAt,
		})
	}
	return resp, nil
}

// CreatePlan создаёт новый учебный план.
func (s *Service) CreatePlan(ctx context.Context, req dto.CreatePlanRequest) (*dto.PlanResponse, error) {
	if req.Title == "" {
		return nil, apperrors.Validationf("title плана не должен быть пустым")
	}
	if req.Period == "" {
		return nil, apperrors.Validationf("period плана не должен быть пустым")
	}
	items := req.Items
	if items == "" {
		items = "[]"
	}
	p := &dbEntities.Plan{
		Title:       req.Title,
		Description: req.Description,
		Period:      req.Period,
		Items:       items,
	}
	if err := s.repo.CreatePlan(ctx, p); err != nil {
		return nil, wrapRepoError("CreatePlan", err)
	}
	return &dto.PlanResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Period:      p.Period,
		Items:       p.Items,
		CreatedAt:   p.CreatedAt,
	}, nil
}

// UpdatePlan обновляет поля учебного плана по ID.
func (s *Service) UpdatePlan(ctx context.Context, id int64, req dto.UpdatePlanRequest) error {
	if id <= 0 {
		return apperrors.Validationf("id плана должен быть больше нуля")
	}
	p := &dbEntities.Plan{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Period:      req.Period,
		Items:       req.Items,
	}
	if err := s.repo.UpdatePlan(ctx, p); err != nil {
		return wrapRepoError("UpdatePlan", err)
	}
	return nil
}

// DeletePlan удаляет учебный план по ID.
func (s *Service) DeletePlan(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperrors.Validationf("id плана должен быть больше нуля")
	}
	if err := s.repo.DeletePlan(ctx, id); err != nil {
		return wrapRepoError("DeletePlan", err)
	}
	return nil
}
