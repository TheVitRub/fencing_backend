// Package dto содержит объекты передачи данных (Data Transfer Objects).
// DTO живут на границе между прикладным слоем и транспортом:
// прикладной слой принимает и возвращает DTO, транспорт маппит их в/из HTTP.
package dto

// LoginRequest — запрос на вход в панель администратора.
type LoginRequest struct {
	Login    string `json:"login"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse — ответ с JWT-токеном при успешном входе.
type LoginResponse struct {
	Token string `json:"token"`
}
