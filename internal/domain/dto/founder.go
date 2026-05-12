package dto

import "time"

// FounderResponse — информация об основателе для отображения на фронте.
type FounderResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	PhotoURL  string    `json:"photo_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpsertFounderRequest — запрос на создание или обновление данных об основателе.
// В таблице всегда одна запись, поэтому операция всегда upsert.
type UpsertFounderRequest struct {
	Name     string `json:"name"     binding:"required"`
	Bio      string `json:"bio"`
	PhotoURL string `json:"photo_url"`
}
