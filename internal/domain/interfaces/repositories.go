// Package interfaces определяет контракты репозиториев.
//
// Прикладной слой (application) работает только с этими интерфейсами,
// не зная о конкретных реализациях. Это позволяет:
//   - подменять реализацию в тестах через моки;
//   - менять БД без правок бизнес-логики.
package interfaces

import (
	"context"

	dbEntities "fencing-club/internal/domain/entities/db"
)

// DBRepository — интерфейс доступа к данным в PostgreSQL.
//
//go:generate go run go.uber.org/mock/mockgen@latest -destination=mocks/db_mock.go -package=mocks fencing-club/internal/domain/interfaces DBRepository
type DBRepository interface {
	// GetAdminByLogin возвращает администратора по логину.
	// Если не найден — возвращает sql.ErrNoRows.
	GetAdminByLogin(ctx context.Context, login string) (*dbEntities.Admin, error)
	GetUserByLogin(ctx context.Context, login string) (*dbEntities.User, error)
	GetUserByID(ctx context.Context, id int64) (*dbEntities.User, error)
	ListUsers(ctx context.Context) ([]dbEntities.User, error)
	CreateUser(ctx context.Context, u *dbEntities.User) error
	UpdateUserRole(ctx context.Context, id int64, role string) error

	// GetPage возвращает страницу по slug.
	// Если не найдена — возвращает sql.ErrNoRows.
	GetPage(ctx context.Context, slug string) (*dbEntities.Page, error)
	// UpsertPage создаёт страницу или обновляет существующую по slug.
	UpsertPage(ctx context.Context, slug, title, content string) error

	// ListEvents возвращает все события, отсортированные по дате убывания.
	ListEvents(ctx context.Context) ([]dbEntities.Event, error)
	GetEvent(ctx context.Context, id int64) (*dbEntities.Event, error)
	// CreateEvent добавляет событие и проставляет ID в переданную структуру.
	CreateEvent(ctx context.Context, e *dbEntities.Event) error
	// UpdateEvent обновляет событие по ID.
	UpdateEvent(ctx context.Context, e *dbEntities.Event) error
	// DeleteEvent удаляет событие по ID.
	DeleteEvent(ctx context.Context, id int64) error

	ListComments(ctx context.Context, targetType string, targetID int64, includeHidden bool) ([]dbEntities.Comment, error)
	CreateComment(ctx context.Context, c *dbEntities.Comment) error
	UpdateCommentStatus(ctx context.Context, id int64, status string) error
	DeleteComment(ctx context.Context, id int64) error

	SetEventAttendance(ctx context.Context, eventID, userID int64, status string) error
	ListEventAttendees(ctx context.Context, eventID int64) ([]dbEntities.EventAttendee, error)
	GetEventAttendance(ctx context.Context, eventID, userID int64) (*dbEntities.EventAttendee, error)
	CreateNotificationsForEventAttendees(ctx context.Context, eventID int64, title, body string) error
	ListNotifications(ctx context.Context, userID int64) ([]dbEntities.Notification, error)
	MarkNotificationRead(ctx context.Context, userID, notificationID int64) error

	ListInstructorProfiles(ctx context.Context) ([]dbEntities.InstructorProfile, error)
	UpsertInstructorProfile(ctx context.Context, p *dbEntities.InstructorProfile) error
	DeleteInstructorProfile(ctx context.Context, userID int64) error

	ListKnowledgeArticles(ctx context.Context, role string) ([]dbEntities.KnowledgeArticle, error)
	CreateKnowledgeArticle(ctx context.Context, a *dbEntities.KnowledgeArticle) error
	UpdateKnowledgeArticle(ctx context.Context, a *dbEntities.KnowledgeArticle) error
	DeleteKnowledgeArticle(ctx context.Context, id int64) error

	ListGlossaryTerms(ctx context.Context) ([]dbEntities.GlossaryTerm, error)
	CreateGlossaryTerm(ctx context.Context, t *dbEntities.GlossaryTerm) error
	UpdateGlossaryTerm(ctx context.Context, t *dbEntities.GlossaryTerm) error
	DeleteGlossaryTerm(ctx context.Context, id int64) error

	ListStudentProgress(ctx context.Context, userID int64) ([]dbEntities.StudentProgress, error)
	UpsertStudentProgress(ctx context.Context, p *dbEntities.StudentProgress) error

	// ListPlans возвращает все учебные планы, отсортированные по дате создания убывания.
	ListPlans(ctx context.Context) ([]dbEntities.Plan, error)
	// CreatePlan добавляет учебный план и проставляет ID в переданную структуру.
	CreatePlan(ctx context.Context, p *dbEntities.Plan) error
	// UpdatePlan обновляет учебный план по ID.
	UpdatePlan(ctx context.Context, p *dbEntities.Plan) error
	// DeletePlan удаляет учебный план по ID.
	DeletePlan(ctx context.Context, id int64) error

	// ListHonorMembers возвращает участников доски почёта, отсортированных по sort_order.
	ListHonorMembers(ctx context.Context) ([]dbEntities.HonorMember, error)
	// CreateHonorMember добавляет участника на доску почёта и проставляет ID.
	CreateHonorMember(ctx context.Context, m *dbEntities.HonorMember) error
	// UpdateHonorMember обновляет участника доски почёта по ID.
	UpdateHonorMember(ctx context.Context, m *dbEntities.HonorMember) error
	// DeleteHonorMember удаляет участника доски почёта по ID.
	DeleteHonorMember(ctx context.Context, id int64) error

	// ListAchievements возвращает достижения школы, отсортированные по году убывания.
	ListAchievements(ctx context.Context) ([]dbEntities.Achievement, error)
	// CreateAchievement добавляет достижение и проставляет ID в переданную структуру.
	CreateAchievement(ctx context.Context, a *dbEntities.Achievement) error
	// UpdateAchievement обновляет достижение по ID.
	UpdateAchievement(ctx context.Context, a *dbEntities.Achievement) error
	// DeleteAchievement удаляет достижение по ID.
	DeleteAchievement(ctx context.Context, id int64) error

	// GetFounder возвращает данные об основателе школы.
	// Если запись ещё не создана — возвращает sql.ErrNoRows.
	GetFounder(ctx context.Context) (*dbEntities.Founder, error)
	// UpsertFounder создаёт или обновляет запись об основателе (id всегда = 1).
	UpsertFounder(ctx context.Context, f *dbEntities.Founder) error
}
