package application

import (
	"context"

	"fencing-club/internal/domain/apperrors"
	dbEntities "fencing-club/internal/domain/entities/db"
	"fencing-club/internal/domain/dto"
)

// ListHonorMembers возвращает участников доски почёта, отсортированных по sort_order.
func (s *Service) ListHonorMembers(ctx context.Context) ([]dto.HonorMemberResponse, error) {
	members, err := s.repo.ListHonorMembers(ctx)
	if err != nil {
		return nil, wrapRepoError("ListHonorMembers", err)
	}
	resp := make([]dto.HonorMemberResponse, 0, len(members))
	for _, m := range members {
		resp = append(resp, dto.HonorMemberResponse{
			ID:          m.ID,
			Name:        m.Name,
			Title:       m.Title,
			Description: m.Description,
			PhotoURL:    m.PhotoURL,
			SortOrder:   m.SortOrder,
			CreatedAt:   m.CreatedAt,
		})
	}
	return resp, nil
}

// CreateHonorMember добавляет участника на доску почёта.
func (s *Service) CreateHonorMember(ctx context.Context, req dto.CreateHonorMemberRequest) (*dto.HonorMemberResponse, error) {
	if req.Name == "" {
		return nil, apperrors.Validationf("name участника не должен быть пустым")
	}
	m := &dbEntities.HonorMember{
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		PhotoURL:    req.PhotoURL,
		SortOrder:   req.SortOrder,
	}
	if err := s.repo.CreateHonorMember(ctx, m); err != nil {
		return nil, wrapRepoError("CreateHonorMember", err)
	}
	return &dto.HonorMemberResponse{
		ID:          m.ID,
		Name:        m.Name,
		Title:       m.Title,
		Description: m.Description,
		PhotoURL:    m.PhotoURL,
		SortOrder:   m.SortOrder,
		CreatedAt:   m.CreatedAt,
	}, nil
}

// UpdateHonorMember обновляет данные участника доски почёта по ID.
func (s *Service) UpdateHonorMember(ctx context.Context, id int64, req dto.UpdateHonorMemberRequest) error {
	if id <= 0 {
		return apperrors.Validationf("id участника должен быть больше нуля")
	}
	m := &dbEntities.HonorMember{
		ID:          id,
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		PhotoURL:    req.PhotoURL,
		SortOrder:   req.SortOrder,
	}
	if err := s.repo.UpdateHonorMember(ctx, m); err != nil {
		return wrapRepoError("UpdateHonorMember", err)
	}
	return nil
}

// DeleteHonorMember удаляет участника доски почёта по ID.
func (s *Service) DeleteHonorMember(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperrors.Validationf("id участника должен быть больше нуля")
	}
	if err := s.repo.DeleteHonorMember(ctx, id); err != nil {
		return wrapRepoError("DeleteHonorMember", err)
	}
	return nil
}
