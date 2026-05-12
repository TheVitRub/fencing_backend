package dto

import "time"

// PlanResponse — учебный план для отображения на фронте.
type PlanResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Period      string    `json:"period"`
	Items       string    `json:"items"` // JSON-массив строк
	CreatedAt   time.Time `json:"created_at"`
}

// CreatePlanRequest — запрос на создание учебного плана.
type CreatePlanRequest struct {
	Title       string `json:"title"  binding:"required"`
	Description string `json:"description"`
	Period      string `json:"period" binding:"required"`
	// Items — JSON-массив строк, например: ["Рапира", "Дага", "Двуручный меч"]
	Items string `json:"items"`
}

// UpdatePlanRequest — запрос на обновление учебного плана.
type UpdatePlanRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Period      string `json:"period"`
	Items       string `json:"items"`
}
