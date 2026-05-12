package http

import (
	"net/http"

	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// listHonorMembersHandler обрабатывает GET /api/honor.
// Возвращает участников доски почёта, отсортированных по sort_order.
func (h *Handler) listHonorMembersHandler(c *gin.Context) {
	members, err := h.svc.ListHonorMembers(c.Request.Context())
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, members)
}

// createHonorMemberHandler обрабатывает POST /api/admin/honor.
// Добавляет участника на доску почёта. Требует JWT.
func (h *Handler) createHonorMemberHandler(c *gin.Context) {
	var req dto.CreateHonorMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	member, err := h.svc.CreateHonorMember(c.Request.Context(), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// updateHonorMemberHandler обрабатывает PUT /api/admin/honor/:id.
// Обновляет данные участника доски почёта по ID. Требует JWT.
func (h *Handler) updateHonorMemberHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}

	var req dto.UpdateHonorMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	if err := h.svc.UpdateHonorMember(c.Request.Context(), id, req); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteHonorMemberHandler обрабатывает DELETE /api/admin/honor/:id.
// Удаляет участника доски почёта по ID. Требует JWT.
func (h *Handler) deleteHonorMemberHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.DeleteHonorMember(c.Request.Context(), id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}
