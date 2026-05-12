package dto

import "time"

// EventResponse — событие для отображения на фронте.
type EventResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateEventRequest — запрос на создание события.
type CreateEventRequest struct {
	Title       string    `json:"title"       binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"        binding:"required"`
	Location    string    `json:"location"`
	ImageURL    string    `json:"image_url"`
}

// UpdateEventRequest — запрос на обновление события.
// ID передаётся в URL-параметре.
type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	ImageURL    string    `json:"image_url"`
}
