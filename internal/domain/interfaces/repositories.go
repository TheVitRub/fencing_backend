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

	// GetPage возвращает страницу по slug.
	// Если не найдена — возвращает sql.ErrNoRows.
	GetPage(ctx context.Context, slug string) (*dbEntities.Page, error)
	// UpsertPage создаёт страницу или обновляет существующую по slug.
	UpsertPage(ctx context.Context, slug, title, content string) error

	// ListEvents возвращает все события, отсортированные по дате убывания.
	ListEvents(ctx context.Context) ([]dbEntities.Event, error)
	// CreateEvent добавляет событие и проставляет ID в переданную структуру.
	CreateEvent(ctx context.Context, e *dbEntities.Event) error
	// UpdateEvent обновляет событие по ID.
	UpdateEvent(ctx context.Context, e *dbEntities.Event) error
	// DeleteEvent удаляет событие по ID.
	DeleteEvent(ctx context.Context, id int64) error

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
