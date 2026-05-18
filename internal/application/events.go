package application

import (
	"context"

	"fencing-club/internal/domain/apperrors"
	"fencing-club/internal/domain/dto"
	dbEntities "fencing-club/internal/domain/entities/db"
)

// ListEvents возвращает все события клуба, отсортированные по дате убывания.
func (s *Service) ListEvents(ctx context.Context) ([]dto.EventResponse, error) {
	events, err := s.repo.ListEvents(ctx)
	if err != nil {
		return nil, wrapRepoError("ListEvents", err)
	}
	resp := make([]dto.EventResponse, 0, len(events))
	for _, e := range events {
		resp = append(resp, dto.EventResponse{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
			Date:        e.Date,
			Location:    e.Location,
			Type:        e.Type,
			Status:      e.Status,
			Discipline:  e.Discipline,
			ImageURL:    e.ImageURL,
			Images:      parseImages(e.Images),
			CreatedAt:   e.CreatedAt,
		})
	}
	return resp, nil
}

// CreateEvent создаёт новое событие и возвращает его с проставленным ID.
func (s *Service) CreateEvent(ctx context.Context, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	if req.Title == "" {
		return nil, apperrors.Validationf("title события не должен быть пустым")
	}
	e := &dbEntities.Event{
		Title:       req.Title,
		Description: req.Description,
		Date:        req.Date,
		Location:    req.Location,
		Type:        defaultString(req.Type, "event"),
		Status:      defaultString(req.Status, "scheduled"),
		Discipline:  defaultString(req.Discipline, "CourtSword"),
		ImageURL:    pickCover(req.ImageURL, req.Images),
		Images:      marshalImages(req.Images),
	}
	if err := s.repo.CreateEvent(ctx, e); err != nil {
		return nil, wrapRepoError("CreateEvent", err)
	}
	s.logger.Info("событие создано",
		s.logger.F("event_id", e.ID),
		s.logger.F("title", e.Title),
	)
	return &dto.EventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Date:        e.Date,
		Location:    e.Location,
		Type:        e.Type,
		Status:      e.Status,
		Discipline:  e.Discipline,
		ImageURL:    e.ImageURL,
		Images:      parseImages(e.Images),
		CreatedAt:   e.CreatedAt,
	}, nil
}

// UpdateEvent обновляет поля события по ID.
func (s *Service) UpdateEvent(ctx context.Context, id int64, req dto.UpdateEventRequest) error {
	if id <= 0 {
		return apperrors.Validationf("id события должен быть больше нуля")
	}
	old, _ := s.repo.GetEvent(ctx, id)
	e := &dbEntities.Event{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Date:        req.Date,
		Location:    req.Location,
		Type:        defaultString(req.Type, "event"),
		Status:      defaultString(req.Status, "scheduled"),
		Discipline:  defaultString(req.Discipline, "CourtSword"),
		ImageURL:    pickCover(req.ImageURL, req.Images),
		Images:      marshalImages(req.Images),
	}
	if err := s.repo.UpdateEvent(ctx, e); err != nil {
		return wrapRepoError("UpdateEvent", err)
	}
	if old != nil && (!old.Date.Equal(e.Date) || old.Location != e.Location || old.Status != e.Status) {
		_ = s.repo.CreateNotificationsForEventAttendees(ctx, id, "Событие изменено", "Проверьте новое время, место или статус: "+e.Title)
	}
	return nil
}

// DeleteEvent удаляет событие по ID.
func (s *Service) DeleteEvent(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperrors.Validationf("id события должен быть больше нуля")
	}
	if err := s.repo.DeleteEvent(ctx, id); err != nil {
		return wrapRepoError("DeleteEvent", err)
	}
	return nil
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
