package dto

import "time"

// AchievementResponse — достижение школы для отображения на фронте.
type AchievementResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Year        int       `json:"year"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateAchievementRequest — запрос на создание достижения.
type CreateAchievementRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Year        int    `json:"year"  binding:"required,min=1900"`
	ImageURL    string `json:"image_url"`
}

// UpdateAchievementRequest — запрос на обновление достижения.
type UpdateAchievementRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        int    `json:"year"`
	ImageURL    string `json:"image_url"`
}
