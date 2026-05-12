package http

import (
	"net/http"

	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// getPageHandler обрабатывает GET /api/pages/:slug.
// Возвращает содержимое страницы по её slug.
func (h *Handler) getPageHandler(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug страницы не указан"})
		return
	}

	page, err := h.svc.GetPage(c.Request.Context(), slug)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, page)
}

// upsertPageHandler обрабатывает PUT /api/admin/pages/:slug.
// Создаёт или обновляет содержимое страницы. Требует JWT.
func (h *Handler) upsertPageHandler(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug страницы не указан"})
		return
	}

	var req dto.UpsertPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	if err := h.svc.UpsertPage(c.Request.Context(), slug, req); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}
