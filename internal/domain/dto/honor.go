package dto

import "time"

// HonorMemberResponse — участник доски почёта для отображения на фронте.
type HonorMemberResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PhotoURL    string    `json:"photo_url"`
	Images      []string  `json:"images"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateHonorMemberRequest — запрос на добавление участника на доску почёта.
type CreateHonorMemberRequest struct {
	Name        string   `json:"name"  binding:"required"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	PhotoURL    string   `json:"photo_url"`
	Images      []string `json:"images"`
	SortOrder   int      `json:"sort_order"`
}

// UpdateHonorMemberRequest — запрос на обновление участника доски почёта.
type UpdateHonorMemberRequest struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	PhotoURL    string   `json:"photo_url"`
	Images      []string `json:"images"`
	SortOrder   int      `json:"sort_order"`
}
