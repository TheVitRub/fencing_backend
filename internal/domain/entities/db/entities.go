// Package db содержит структуры-сущности, которые точно отражают схему БД.
// Теги db:"..." используются sqlx для автоматического маппинга строк.
// Эти структуры НЕ должны утекать в транспортный слой — для этого есть dto.
package db

import "time"

// Admin — запись администратора сайта.
type Admin struct {
	ID           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
}

// User — пользователь сайта: зарегистрированный гость, ученик, инструктор, админ или основатель.
type User struct {
	ID           int64     `db:"id"`
	Login        string    `db:"login"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	DisplayName  string    `db:"display_name"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Page — страница сайта с произвольным JSON-контентом.
// Поле Content хранит JSON-строку, которую фронт разбирает самостоятельно.
type Page struct {
	ID        int64     `db:"id"`
	Slug      string    `db:"slug"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Event — событие (турнир, тренировка, показательное выступление).
// Поле Images хранит JSON-массив URL изображений (галерея).
type Event struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Date        time.Time `db:"date"`
	Location    string    `db:"location"`
	Type        string    `db:"type"`
	Status      string    `db:"status"`
	Discipline  string    `db:"discipline"`
	ImageURL    string    `db:"image_url"` // главная обложка (deprecated, оставлено для совместимости)
	Images      string    `db:"images"`    // JSON-массив URL: ["/uploads/a.jpg", ...]
	CreatedAt   time.Time `db:"created_at"`
}

// Plan — учебный план на определённый период.
// Поле Items хранит JSON-массив строк с пунктами плана.
type Plan struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Period      string    `db:"period"`
	Items       string    `db:"items"`
	CreatedAt   time.Time `db:"created_at"`
}

// HonorMember — участник доски почёта клуба.
type HonorMember struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	PhotoURL    string    `db:"photo_url"`
	Images      string    `db:"images"` // JSON-массив URL дополнительных фото
	SortOrder   int       `db:"sort_order"`
	CreatedAt   time.Time `db:"created_at"`
}

// Achievement — достижение школы (победа на турнире, награда и т.д.).
type Achievement struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Year        int       `db:"year"`
	ImageURL    string    `db:"image_url"` // главная обложка
	Images      string    `db:"images"`    // JSON-массив URL дополнительных изображений
	CreatedAt   time.Time `db:"created_at"`
}

// Founder — основатель школы. В таблице всегда одна запись с id=1.
type Founder struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Bio       string    `db:"bio"`
	PhotoURL  string    `db:"photo_url"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Comment struct {
	ID              int64     `db:"id"`
	TargetType      string    `db:"target_type"`
	TargetID        int64     `db:"target_id"`
	UserID          int64     `db:"user_id"`
	UserDisplayName string    `db:"user_display_name"`
	UserRole        string    `db:"user_role"`
	Body            string    `db:"body"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type EventAttendee struct {
	EventID         int64     `db:"event_id"`
	UserID          int64     `db:"user_id"`
	UserDisplayName string    `db:"user_display_name"`
	UserRole        string    `db:"user_role"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
}

type Notification struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	EventID   *int64    `db:"event_id"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	IsRead    bool      `db:"is_read"`
	CreatedAt time.Time `db:"created_at"`
}

type InstructorProfile struct {
	ID             int64     `db:"id"`
	UserID         int64     `db:"user_id"`
	Role           string    `db:"role"`
	Name           string    `db:"name"`
	PhotoURL       string    `db:"photo_url"`
	Specialization string    `db:"specialization"`
	Weapons        string    `db:"weapons"`
	Experience     string    `db:"experience"`
	Quote          string    `db:"quote"`
	Bio            string    `db:"bio"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type KnowledgeArticle struct {
	ID         int64     `db:"id"`
	Title      string    `db:"title"`
	Category   string    `db:"category"`
	Body       string    `db:"body"`
	Visibility string    `db:"visibility"`
	SortOrder  int       `db:"sort_order"`
	ImageURL   string    `db:"image_url"`
	Images     string    `db:"images"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type GlossaryTerm struct {
	ID         int64     `db:"id"`
	Term       string    `db:"term"`
	Category   string    `db:"category"`
	Definition string    `db:"definition"`
	ImageURL   string    `db:"image_url"`
	Images     string    `db:"images"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type StudentProgress struct {
	ID              int64     `db:"id"`
	UserID          int64     `db:"user_id"`
	UserDisplayName string    `db:"user_display_name"`
	Discipline      string    `db:"discipline"`
	Level           string    `db:"level"`
	PassedChecks    string    `db:"passed_checks"`
	InstructorNote  string    `db:"instructor_note"`
	UpdatedBy       *int64    `db:"updated_by"`
	UpdatedAt       time.Time `db:"updated_at"`
}
