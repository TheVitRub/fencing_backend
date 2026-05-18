package dto

import "time"

// EventResponse — событие для отображения на фронте.
type EventResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	ImageURL    string    `json:"image_url"` // главная обложка
	Images      []string  `json:"images"`    // галерея
	CreatedAt   time.Time `json:"created_at"`
}

// CreateEventRequest — запрос на создание события.
type CreateEventRequest struct {
	Title       string    `json:"title"       binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"        binding:"required"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	ImageURL    string    `json:"image_url"`
	Images      []string  `json:"images"`
}

// UpdateEventRequest — запрос на обновление события.
// ID передаётся в URL-параметре.
type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	ImageURL    string    `json:"image_url"`
	Images      []string  `json:"images"`
}
