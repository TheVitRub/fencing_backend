package dto

import "time"

type CommentResponse struct {
	ID              int64     `json:"id"`
	TargetType      string    `json:"target_type"`
	TargetID        int64     `json:"target_id"`
	UserID          int64     `json:"user_id"`
	UserDisplayName string    `json:"user_display_name"`
	UserRole        string    `json:"user_role"`
	Body            string    `json:"body"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   int64  `json:"target_id" binding:"required"`
	Body       string `json:"body" binding:"required"`
}

type UpdateCommentStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type EventAttendeeResponse struct {
	EventID         int64     `json:"event_id"`
	UserID          int64     `json:"user_id"`
	UserDisplayName string    `json:"user_display_name"`
	UserRole        string    `json:"user_role"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type NotificationResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	EventID   *int64    `json:"event_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type InstructorProfileResponse struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Role           string    `json:"role"`
	Name           string    `json:"name"`
	PhotoURL       string    `json:"photo_url"`
	Specialization string    `json:"specialization"`
	Weapons        string    `json:"weapons"`
	Experience     string    `json:"experience"`
	Quote          string    `json:"quote"`
	Bio            string    `json:"bio"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UpsertInstructorProfileRequest struct {
	UserID         int64  `json:"user_id"`
	Name           string `json:"name"`
	PhotoURL       string `json:"photo_url"`
	Specialization string `json:"specialization"`
	Weapons        string `json:"weapons"`
	Experience     string `json:"experience"`
	Quote          string `json:"quote"`
	Bio            string `json:"bio"`
}

type KnowledgeArticleResponse struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	Body       string    `json:"body"`
	Visibility string    `json:"visibility"`
	SortOrder  int       `json:"sort_order"`
	ImageURL   string    `json:"image_url"`
	Images     []string  `json:"images"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpsertKnowledgeArticleRequest struct {
	Title      string   `json:"title" binding:"required"`
	Category   string   `json:"category"`
	Body       string   `json:"body"`
	Visibility string   `json:"visibility"`
	SortOrder  int      `json:"sort_order"`
	ImageURL   string   `json:"image_url"`
	Images     []string `json:"images"`
}

type GlossaryTermResponse struct {
	ID         int64     `json:"id"`
	Term       string    `json:"term"`
	Category   string    `json:"category"`
	Definition string    `json:"definition"`
	ImageURL   string    `json:"image_url"`
	Images     []string  `json:"images"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpsertGlossaryTermRequest struct {
	Term       string   `json:"term" binding:"required"`
	Category   string   `json:"category"`
	Definition string   `json:"definition"`
	ImageURL   string   `json:"image_url"`
	Images     []string `json:"images"`
}

type StudentProgressResponse struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	UserDisplayName string    `json:"user_display_name"`
	Discipline      string    `json:"discipline"`
	Level           string    `json:"level"`
	PassedChecks    []string  `json:"passed_checks"`
	InstructorNote  string    `json:"instructor_note"`
	UpdatedBy       *int64    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UpsertStudentProgressRequest struct {
	UserID         int64    `json:"user_id" binding:"required"`
	Discipline     string   `json:"discipline" binding:"required"`
	Level          string   `json:"level"`
	PassedChecks   []string `json:"passed_checks"`
	InstructorNote string   `json:"instructor_note"`
}
