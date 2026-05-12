package http

import (
	"net/http"
	"strconv"

	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// listEventsHandler обрабатывает GET /api/events.
// Возвращает все события, отсортированные по дате убывания.
func (h *Handler) listEventsHandler(c *gin.Context) {
	events, err := h.svc.ListEvents(c.Request.Context())
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, events)
}

// createEventHandler обрабатывает POST /api/admin/events.
// Создаёт новое событие и возвращает его с проставленным ID. Требует JWT.
func (h *Handler) createEventHandler(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	event, err := h.svc.CreateEvent(c.Request.Context(), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// updateEventHandler обрабатывает PUT /api/admin/events/:id.
// Обновляет поля события по ID. Требует JWT.
func (h *Handler) updateEventHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	if err := h.svc.UpdateEvent(c.Request.Context(), id, req); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteEventHandler обрабатывает DELETE /api/admin/events/:id.
// Удаляет событие по ID. Требует JWT.
func (h *Handler) deleteEventHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}

	if err := h.svc.DeleteEvent(c.Request.Context(), id); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

// parseIDParam извлекает числовой параметр из URL и записывает 400 при ошибке.
func parseIDParam(c *gin.Context, name string) (int64, error) {
	raw := c.Param(name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр " + name + " должен быть положительным числом"})
		return 0, err
	}
	return id, nil
}
