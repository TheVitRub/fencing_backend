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
