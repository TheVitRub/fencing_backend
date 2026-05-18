package application

import (
	"context"
	"strings"

	"fencing-club/internal/domain/apperrors"
	"fencing-club/internal/domain/dto"
	dbEntities "fencing-club/internal/domain/entities/db"
)

func commentToDTO(c dbEntities.Comment) dto.CommentResponse {
	return dto.CommentResponse{
		ID:              c.ID,
		TargetType:      c.TargetType,
		TargetID:        c.TargetID,
		UserID:          c.UserID,
		UserDisplayName: c.UserDisplayName,
		UserRole:        c.UserRole,
		Body:            c.Body,
		Status:          c.Status,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

func (s *Service) ListComments(ctx context.Context, targetType string, targetID int64, role string) ([]dto.CommentResponse, error) {
	comments, err := s.repo.ListComments(ctx, targetType, targetID, IsInstructorRole(role))
	if err != nil {
		return nil, wrapRepoError("ListComments", err)
	}
	resp := make([]dto.CommentResponse, 0, len(comments))
	for _, c := range comments {
		resp = append(resp, commentToDTO(c))
	}
	return resp, nil
}

func (s *Service) CreateComment(ctx context.Context, userID int64, req dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	if userID <= 0 {
		return nil, apperrors.Unauthorizedf("нужен вход")
	}
	if req.TargetID <= 0 || !validCommentTarget(req.TargetType) {
		return nil, apperrors.Validationf("некорректная цель комментария")
	}
	body := strings.TrimSpace(req.Body)
	if body == "" {
		return nil, apperrors.Validationf("комментарий не должен быть пустым")
	}
	c := &dbEntities.Comment{TargetType: req.TargetType, TargetID: req.TargetID, UserID: userID, Body: body}
	if err := s.repo.CreateComment(ctx, c); err != nil {
		return nil, wrapRepoError("CreateComment", err)
	}
	comments, _ := s.repo.ListComments(ctx, req.TargetType, req.TargetID, true)
	for _, comment := range comments {
		if comment.ID == c.ID {
			resp := commentToDTO(comment)
			return &resp, nil
		}
	}
	return &dto.CommentResponse{ID: c.ID, TargetType: req.TargetType, TargetID: req.TargetID, UserID: userID, Body: body, Status: "visible"}, nil
}

func (s *Service) UpdateCommentStatus(ctx context.Context, id int64, status string) error {
	if status != "visible" && status != "hidden" && status != "deleted" {
		return apperrors.Validationf("неизвестный статус комментария")
	}
	return wrapRepoError("UpdateCommentStatus", s.repo.UpdateCommentStatus(ctx, id, status))
}

func validCommentTarget(value string) bool {
	return value == "event" || value == "achievement" || value == "plan"
}

func (s *Service) SetEventAttendance(ctx context.Context, eventID, userID int64, status string) error {
	if eventID <= 0 || userID <= 0 {
		return apperrors.Validationf("некорректные параметры записи")
	}
	if status != "going" && status != "cancelled" {
		status = "going"
	}
	return wrapRepoError("SetEventAttendance", s.repo.SetEventAttendance(ctx, eventID, userID, status))
}

func (s *Service) ListEventAttendees(ctx context.Context, eventID int64) ([]dto.EventAttendeeResponse, error) {
	items, err := s.repo.ListEventAttendees(ctx, eventID)
	if err != nil {
		return nil, wrapRepoError("ListEventAttendees", err)
	}
	resp := make([]dto.EventAttendeeResponse, 0, len(items))
	for _, a := range items {
		resp = append(resp, dto.EventAttendeeResponse{
			EventID: a.EventID, UserID: a.UserID, UserDisplayName: a.UserDisplayName,
			UserRole: a.UserRole, Status: a.Status, CreatedAt: a.CreatedAt,
		})
	}
	return resp, nil
}

func (s *Service) ListNotifications(ctx context.Context, userID int64) ([]dto.NotificationResponse, error) {
	items, err := s.repo.ListNotifications(ctx, userID)
	if err != nil {
		return nil, wrapRepoError("ListNotifications", err)
	}
	resp := make([]dto.NotificationResponse, 0, len(items))
	for _, n := range items {
		resp = append(resp, dto.NotificationResponse(n))
	}
	return resp, nil
}

func (s *Service) MarkNotificationRead(ctx context.Context, userID, notificationID int64) error {
	return wrapRepoError("MarkNotificationRead", s.repo.MarkNotificationRead(ctx, userID, notificationID))
}

func (s *Service) ListInstructorProfiles(ctx context.Context) ([]dto.InstructorProfileResponse, error) {
	items, err := s.repo.ListInstructorProfiles(ctx)
	if err != nil {
		return nil, wrapRepoError("ListInstructorProfiles", err)
	}
	resp := make([]dto.InstructorProfileResponse, 0, len(items))
	for _, p := range items {
		resp = append(resp, dto.InstructorProfileResponse(p))
	}
	return resp, nil
}

func (s *Service) UpsertInstructorProfile(ctx context.Context, actorID int64, actorRole string, req dto.UpsertInstructorProfileRequest) error {
	userID := req.UserID
	if userID == 0 {
		userID = actorID
	}
	if userID != actorID && !IsAdminRole(actorRole) {
		return apperrors.Forbiddenf("нельзя редактировать чужой профиль")
	}
	p := &dbEntities.InstructorProfile{
		UserID: userID, Name: req.Name, PhotoURL: req.PhotoURL, Specialization: req.Specialization,
		Weapons: req.Weapons, Experience: req.Experience, Quote: req.Quote, Bio: req.Bio,
	}
	return wrapRepoError("UpsertInstructorProfile", s.repo.UpsertInstructorProfile(ctx, p))
}

func (s *Service) ListKnowledgeArticles(ctx context.Context, role string) ([]dto.KnowledgeArticleResponse, error) {
	items, err := s.repo.ListKnowledgeArticles(ctx, role)
	if err != nil {
		return nil, wrapRepoError("ListKnowledgeArticles", err)
	}
	resp := make([]dto.KnowledgeArticleResponse, 0, len(items))
	for _, a := range items {
		resp = append(resp, dto.KnowledgeArticleResponse(a))
	}
	return resp, nil
}

func (s *Service) CreateKnowledgeArticle(ctx context.Context, req dto.UpsertKnowledgeArticleRequest) (*dto.KnowledgeArticleResponse, error) {
	a := &dbEntities.KnowledgeArticle{Title: req.Title, Category: req.Category, Body: req.Body, Visibility: defaultString(req.Visibility, "public"), SortOrder: req.SortOrder}
	if err := s.repo.CreateKnowledgeArticle(ctx, a); err != nil {
		return nil, wrapRepoError("CreateKnowledgeArticle", err)
	}
	return &dto.KnowledgeArticleResponse{ID: a.ID, Title: a.Title, Category: a.Category, Body: a.Body, Visibility: a.Visibility, SortOrder: a.SortOrder}, nil
}

func (s *Service) UpdateKnowledgeArticle(ctx context.Context, id int64, req dto.UpsertKnowledgeArticleRequest) error {
	a := &dbEntities.KnowledgeArticle{ID: id, Title: req.Title, Category: req.Category, Body: req.Body, Visibility: defaultString(req.Visibility, "public"), SortOrder: req.SortOrder}
	return wrapRepoError("UpdateKnowledgeArticle", s.repo.UpdateKnowledgeArticle(ctx, a))
}

func (s *Service) DeleteKnowledgeArticle(ctx context.Context, id int64) error {
	return wrapRepoError("DeleteKnowledgeArticle", s.repo.DeleteKnowledgeArticle(ctx, id))
}

func (s *Service) ListGlossaryTerms(ctx context.Context) ([]dto.GlossaryTermResponse, error) {
	items, err := s.repo.ListGlossaryTerms(ctx)
	if err != nil {
		return nil, wrapRepoError("ListGlossaryTerms", err)
	}
	resp := make([]dto.GlossaryTermResponse, 0, len(items))
	for _, t := range items {
		resp = append(resp, dto.GlossaryTermResponse(t))
	}
	return resp, nil
}

func (s *Service) CreateGlossaryTerm(ctx context.Context, req dto.UpsertGlossaryTermRequest) (*dto.GlossaryTermResponse, error) {
	t := &dbEntities.GlossaryTerm{Term: req.Term, Category: req.Category, Definition: req.Definition}
	if err := s.repo.CreateGlossaryTerm(ctx, t); err != nil {
		return nil, wrapRepoError("CreateGlossaryTerm", err)
	}
	return &dto.GlossaryTermResponse{ID: t.ID, Term: t.Term, Category: t.Category, Definition: t.Definition}, nil
}

func (s *Service) UpdateGlossaryTerm(ctx context.Context, id int64, req dto.UpsertGlossaryTermRequest) error {
	t := &dbEntities.GlossaryTerm{ID: id, Term: req.Term, Category: req.Category, Definition: req.Definition}
	return wrapRepoError("UpdateGlossaryTerm", s.repo.UpdateGlossaryTerm(ctx, t))
}

func (s *Service) DeleteGlossaryTerm(ctx context.Context, id int64) error {
	return wrapRepoError("DeleteGlossaryTerm", s.repo.DeleteGlossaryTerm(ctx, id))
}

func (s *Service) ListStudentProgress(ctx context.Context, userID int64) ([]dto.StudentProgressResponse, error) {
	items, err := s.repo.ListStudentProgress(ctx, userID)
	if err != nil {
		return nil, wrapRepoError("ListStudentProgress", err)
	}
	resp := make([]dto.StudentProgressResponse, 0, len(items))
	for _, p := range items {
		resp = append(resp, dto.StudentProgressResponse{
			ID: p.ID, UserID: p.UserID, UserDisplayName: p.UserDisplayName, Discipline: p.Discipline,
			Level: p.Level, PassedChecks: parseImages(p.PassedChecks), InstructorNote: p.InstructorNote,
			UpdatedBy: p.UpdatedBy, UpdatedAt: p.UpdatedAt,
		})
	}
	return resp, nil
}

func (s *Service) UpsertStudentProgress(ctx context.Context, actorID int64, req dto.UpsertStudentProgressRequest) error {
	p := &dbEntities.StudentProgress{
		UserID: req.UserID, Discipline: req.Discipline, Level: req.Level, PassedChecks: marshalImages(req.PassedChecks),
		InstructorNote: req.InstructorNote, UpdatedBy: &actorID,
	}
	return wrapRepoError("UpsertStudentProgress", s.repo.UpsertStudentProgress(ctx, p))
}
