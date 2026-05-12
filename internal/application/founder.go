package application

import (
	"context"

	"fencing-club/internal/domain/apperrors"
	dbEntities "fencing-club/internal/domain/entities/db"
	"fencing-club/internal/domain/dto"
)

// GetFounder возвращает информацию об основателе школы.
// Если запись ещё не создана — возвращает apperrors.KindNotFound.
func (s *Service) GetFounder(ctx context.Context) (*dto.FounderResponse, error) {
	f, err := s.repo.GetFounder(ctx)
	if err != nil {
		return nil, wrapRepoError("GetFounder", err)
	}
	return &dto.FounderResponse{
		ID:        f.ID,
		Name:      f.Name,
		Bio:       f.Bio,
		PhotoURL:  f.PhotoURL,
		UpdatedAt: f.UpdatedAt,
	}, nil
}

// UpsertFounder создаёт или обновляет запись об основателе.
func (s *Service) UpsertFounder(ctx context.Context, req dto.UpsertFounderRequest) error {
	if req.Name == "" {
		return apperrors.Validationf("name основателя не должен быть пустым")
	}
	f := &dbEntities.Founder{
		Name:     req.Name,
		Bio:      req.Bio,
		PhotoURL: req.PhotoURL,
	}
	if err := s.repo.UpsertFounder(ctx, f); err != nil {
		return wrapRepoError("UpsertFounder", err)
	}
	return nil
}
