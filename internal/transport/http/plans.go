package http

import (
	"net/http"

	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// listPlansHandler обрабатывает GET /api/plans.
// Возвращает все учебные планы.
func (h *Handler) listPlansHandler(c *gin.Context) {
	plans, err := h.svc.ListPlans(c.Request.Context())
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, plans)
}

// createPlanHandler обрабатывает POST /api/admin/plans.
// Создаёт новый учебный план. Требует JWT.
func (h *Handler) createPlanHandler(c *gin.Context) {
	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	plan, err := h.svc.CreatePlan(c.Request.Context(), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// updatePlanHandler обрабатывает PUT /api/admin/plans/:id.
// Обновляет учебный план по ID. Требует JWT.
func (h *Handler) updatePlanHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}

	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	if err := h.svc.UpdatePlan(c.Request.Context(), id, req); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

// deletePlanHandler обрабатывает DELETE /api/admin/plans/:id.
// Удаляет учебный план по ID. Требует JWT.
func (h *Handler) deletePlanHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.DeletePlan(c.Request.Context(), id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}
