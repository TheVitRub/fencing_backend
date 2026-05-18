package http

import (
	"net/http"

	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// loginHandler обрабатывает POST /api/auth/login.
// Принимает логин/пароль, возвращает JWT-токен при успехе.
func (h *Handler) loginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) registerHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	resp, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) meHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	resp, err := h.svc.Me(c.Request.Context(), claims.UserID)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, resp)
}
