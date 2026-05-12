package http

import (
	"net/http"

	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// listAchievementsHandler обрабатывает GET /api/achievements.
// Возвращает достижения школы, отсортированные по году убывания.
func (h *Handler) listAchievementsHandler(c *gin.Context) {
	achievements, err := h.svc.ListAchievements(c.Request.Context())
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, achievements)
}

// createAchievementHandler обрабатывает POST /api/admin/achievements.
// Добавляет достижение школы. Требует JWT.
func (h *Handler) createAchievementHandler(c *gin.Context) {
	var req dto.CreateAchievementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	achievement, err := h.svc.CreateAchievement(c.Request.Context(), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, achievement)
}

// updateAchievementHandler обрабатывает PUT /api/admin/achievements/:id.
// Обновляет достижение по ID. Требует JWT.
func (h *Handler) updateAchievementHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}

	var req dto.UpdateAchievementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	if err := h.svc.UpdateAchievement(c.Request.Context(), id, req); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteAchievementHandler обрабатывает DELETE /api/admin/achievements/:id.
// Удаляет достижение по ID. Требует JWT.
func (h *Handler) deleteAchievementHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.DeleteAchievement(c.Request.Context(), id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}
