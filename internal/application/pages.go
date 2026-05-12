package application

import (
	"context"

	"fencing-club/internal/domain/dto"
)

// GetPage возвращает страницу по её slug.
// Если страница не найдена — возвращает apperrors.KindNotFound.
func (s *Service) GetPage(ctx context.Context, slug string) (*dto.PageResponse, error) {
	page, err := s.repo.GetPage(ctx, slug)
	if err != nil {
		return nil, wrapRepoError("GetPage slug="+slug, err)
	}
	return &dto.PageResponse{
		ID:        page.ID,
		Slug:      page.Slug,
		Title:     page.Title,
		Content:   page.Content,
		UpdatedAt: page.UpdatedAt,
	}, nil
}

// UpsertPage создаёт или обновляет страницу с указанным slug.
func (s *Service) UpsertPage(ctx context.Context, slug string, req dto.UpsertPageRequest) error {
	if err := s.repo.UpsertPage(ctx, slug, req.Title, req.Content); err != nil {
		return wrapRepoError("UpsertPage slug="+slug, err)
	}
	s.logger.Info("страница обновлена", s.logger.F("slug", slug))
	return nil
}
