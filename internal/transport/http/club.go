package http

import (
	"net/http"
	"strconv"

	"fencing-club/internal/application"
	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) listUsersHandler(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) listStudentUsersHandler(c *gin.Context) {
	users, err := h.svc.ListUsersByRoles(c.Request.Context(), "student")
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) listInstructorUsersHandler(c *gin.Context) {
	users, err := h.svc.ListUsersByRoles(c.Request.Context(), "instructor", "admin", "founder")
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) updateUserRoleHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	if err := h.svc.UpdateUserRole(c.Request.Context(), id, req.Role); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) listCommentsHandler(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID, _ := strconv.ParseInt(c.Query("target_id"), 10, 64)
	role := ""
	if claims, ok := getClaimsOptional(c); ok {
		role = claims.Role
	}
	comments, err := h.svc.ListComments(c.Request.Context(), targetType, targetID, role)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *Handler) createCommentHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	comment, err := h.svc.CreateComment(c.Request.Context(), claims.UserID, req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) updateCommentStatusHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	var req dto.UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	if err := h.svc.UpdateCommentStatus(c.Request.Context(), id, req.Status); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) setAttendanceHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	eventID, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	status := c.DefaultQuery("status", "going")
	if err := h.svc.SetEventAttendance(c.Request.Context(), eventID, claims.UserID, status); err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) listAttendeesHandler(c *gin.Context) {
	eventID, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	items, svcErr := h.svc.ListEventAttendees(c.Request.Context(), eventID)
	if svcErr != nil {
		st, msg := mapError(svcErr)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) listNotificationsHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	items, err := h.svc.ListNotifications(c.Request.Context(), claims.UserID)
	if err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) markNotificationReadHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	if err := h.svc.MarkNotificationRead(c.Request.Context(), claims.UserID, id); err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) listInstructorProfilesHandler(c *gin.Context) {
	items, err := h.svc.ListInstructorProfiles(c.Request.Context())
	if err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) upsertInstructorProfileHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	var req dto.UpsertInstructorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	if err := h.svc.UpsertInstructorProfile(c.Request.Context(), claims.UserID, claims.Role, req); err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) listKnowledgeHandler(c *gin.Context) {
	role := ""
	if claims, ok := getClaimsOptional(c); ok {
		role = claims.Role
	}
	items, err := h.svc.ListKnowledgeArticles(c.Request.Context(), role)
	if err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) createKnowledgeHandler(c *gin.Context) {
	var req dto.UpsertKnowledgeArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	item, err := h.svc.CreateKnowledgeArticle(c.Request.Context(), req)
	if err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) updateKnowledgeHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	var req dto.UpsertKnowledgeArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	if err := h.svc.UpdateKnowledgeArticle(c.Request.Context(), id, req); err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) deleteKnowledgeHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	if err := h.svc.DeleteKnowledgeArticle(c.Request.Context(), id); err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) listGlossaryHandler(c *gin.Context) {
	items, err := h.svc.ListGlossaryTerms(c.Request.Context())
	if err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) createGlossaryHandler(c *gin.Context) {
	var req dto.UpsertGlossaryTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	item, err := h.svc.CreateGlossaryTerm(c.Request.Context(), req)
	if err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) updateGlossaryHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	var req dto.UpsertGlossaryTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	if err := h.svc.UpdateGlossaryTerm(c.Request.Context(), id, req); err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) deleteGlossaryHandler(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return
	}
	if err := h.svc.DeleteGlossaryTerm(c.Request.Context(), id); err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) listProgressHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	userID := claims.UserID
	if claims.Role == "instructor" || claims.Role == "admin" || claims.Role == "founder" {
		userID = 0
	}
	items, err := h.svc.ListStudentProgress(c.Request.Context(), userID)
	if err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) upsertProgressHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	var req dto.UpsertStudentProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	if err := h.svc.UpsertStudentProgress(c.Request.Context(), claims.UserID, req); err != nil {
		st, msg := mapError(err)
		c.JSON(st, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func getClaimsOptional(c *gin.Context) (*application.Claims, bool) {
	raw, ok := c.Get("claims")
	if !ok {
		return nil, false
	}
	if claims, ok := raw.(*application.Claims); ok {
		return claims, true
	}
	return nil, false
}
