package http

import (
	"net/http"

	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// getFounderHandler обрабатывает GET /api/founder.
// Возвращает информацию об основателе школы.
func (h *Handler) getFounderHandler(c *gin.Context) {
	founder, err := h.svc.GetFounder(c.Request.Context())
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, founder)
}

// upsertFounderHandler обрабатывает PUT /api/admin/founder.
// Создаёт или обновляет запись об основателе. Требует JWT.
func (h *Handler) upsertFounderHandler(c *gin.Context) {
	var req dto.UpsertFounderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	if err := h.svc.UpsertFounder(c.Request.Context(), req); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}
