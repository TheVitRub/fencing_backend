package dto

import "time"

// PageResponse — страница для отображения на фронте.
type PageResponse struct {
	ID        int64     `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Content   string    `json:"content"` // JSON-строка с произвольным контентом
	UpdatedAt time.Time `json:"updated_at"`
}

// UpsertPageRequest — запрос на создание или обновление страницы.
// Slug передаётся в URL, поэтому здесь его нет.
type UpsertPageRequest struct {
	Title   string `json:"title"   binding:"required"`
	Content string `json:"content" binding:"required"`
}
