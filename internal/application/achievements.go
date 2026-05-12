package application

import (
	"context"

	"fencing-club/internal/domain/apperrors"
	dbEntities "fencing-club/internal/domain/entities/db"
	"fencing-club/internal/domain/dto"
)

// ListAchievements возвращает достижения школы, отсортированные по году убывания.
func (s *Service) ListAchievements(ctx context.Context) ([]dto.AchievementResponse, error) {
	list, err := s.repo.ListAchievements(ctx)
	if err != nil {
		return nil, wrapRepoError("ListAchievements", err)
	}
	resp := make([]dto.AchievementResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, dto.AchievementResponse{
			ID:          a.ID,
			Title:       a.Title,
			Description: a.Description,
			Year:        a.Year,
			ImageURL:    a.ImageURL,
			CreatedAt:   a.CreatedAt,
		})
	}
	return resp, nil
}

// CreateAchievement добавляет достижение школы.
func (s *Service) CreateAchievement(ctx context.Context, req dto.CreateAchievementRequest) (*dto.AchievementResponse, error) {
	if req.Title == "" {
		return nil, apperrors.Validationf("title достижения не должен быть пустым")
	}
	if req.Year < 1900 {
		return nil, apperrors.Validationf("year достижения должен быть не меньше 1900")
	}
	a := &dbEntities.Achievement{
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		ImageURL:    req.ImageURL,
	}
	if err := s.repo.CreateAchievement(ctx, a); err != nil {
		return nil, wrapRepoError("CreateAchievement", err)
	}
	return &dto.AchievementResponse{
		ID:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		Year:        a.Year,
		ImageURL:    a.ImageURL,
		CreatedAt:   a.CreatedAt,
	}, nil
}

// UpdateAchievement обновляет поля достижения по ID.
func (s *Service) UpdateAchievement(ctx context.Context, id int64, req dto.UpdateAchievementRequest) error {
	if id <= 0 {
		return apperrors.Validationf("id достижения должен быть больше нуля")
	}
	a := &dbEntities.Achievement{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		ImageURL:    req.ImageURL,
	}
	if err := s.repo.UpdateAchievement(ctx, a); err != nil {
		return wrapRepoError("UpdateAchievement", err)
	}
	return nil
}

// DeleteAchievement удаляет достижение по ID.
func (s *Service) DeleteAchievement(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperrors.Validationf("id достижения должен быть больше нуля")
	}
	if err := s.repo.DeleteAchievement(ctx, id); err != nil {
		return wrapRepoError("DeleteAchievement", err)
	}
	return nil
}
