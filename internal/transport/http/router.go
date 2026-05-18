package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fencing-club/internal/application"
	"fencing-club/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Handler — корневой объект HTTP-слоя.
// Содержит ссылку на сервис приложения, логгер и параметры загрузки файлов.
type Handler struct {
	svc    *application.Service
	logger *logger.Logger
	upload UploadConfig
}

// NewRouter создаёт Gin-роутер с зарегистрированными маршрутами.
// Публичные маршруты доступны без токена; маршруты /admin требуют JWT.
// Статическая раздача файлов из upload.Dir подвешена на upload.URLPrefix.
func NewRouter(svc *application.Service, log *logger.Logger, upload UploadConfig) *gin.Engine {
	h := &Handler{svc: svc, logger: log, upload: upload}

	// Создаём директорию загрузок, если её нет — иначе SaveUploadedFile упадёт.
	if upload.Dir != "" {
		abs, _ := filepath.Abs(upload.Dir)
		if err := os.MkdirAll(upload.Dir, 0o755); err != nil {
			log.Warn("не удалось создать директорию загрузок",
				log.F("dir", abs),
				log.F("error", err),
			)
		} else {
			log.Info("директория загрузок готова", log.F("dir", abs))
		}
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// ── Статическая раздача загруженных файлов ──────────────────────────────
	// URL: /uploads/abc123.jpg  →  файл из <upload.Dir>/abc123.jpg
	if upload.Dir != "" && upload.URLPrefix != "" {
		prefix := "/" + strings.Trim(upload.URLPrefix, "/")
		r.Static(prefix, upload.Dir)
	}

	api := r.Group("/api")
	{
		// --- Аутентификация ---
		auth := api.Group("/auth")
		auth.POST("/login", h.loginHandler)
		auth.POST("/register", h.registerHandler)
		auth.GET("/oauth/:provider/start", h.oauthStartHandler)
		auth.GET("/oauth/:provider/callback", h.oauthCallbackHandler)

		// --- Публичные маршруты ---
		api.GET("/pages/:slug", h.getPageHandler)
		api.GET("/events", h.listEventsHandler)
		api.GET("/plans", h.listPlansHandler)
		api.GET("/honor", h.listHonorMembersHandler)
		api.GET("/achievements", h.listAchievementsHandler)
		api.GET("/founder", h.getFounderHandler)
		api.GET("/comments", h.listCommentsHandler)
		api.GET("/events/:id/attendees", h.listAttendeesHandler)
		api.GET("/instructors", h.listInstructorProfilesHandler)
		api.GET("/knowledge", h.listKnowledgeHandler)
		api.GET("/glossary", h.listGlossaryHandler)

		// --- Административные маршруты (требуют JWT) ---
		admin := api.Group("/admin")
		admin.Use(authMiddleware(svc))
		{
			admin.GET("/me", h.meHandler)
			admin.GET("/knowledge-view", h.listKnowledgeHandler)
			admin.GET("/notifications", h.listNotificationsHandler)
			admin.PUT("/notifications/:id/read", h.markNotificationReadHandler)
			admin.GET("/comments", h.listCommentsHandler)
			admin.POST("/comments", h.createCommentHandler)
			admin.PUT("/comments/:id/status", roleMiddleware("instructor", "admin", "founder"), h.updateCommentStatusHandler)
			admin.PUT("/events/:id/attendance", h.setAttendanceHandler)
			admin.GET("/progress", h.listProgressHandler)

			// Загрузка файлов
			admin.POST("/upload", h.uploadHandler)

			// Страницы
			admin.PUT("/pages/:slug", h.upsertPageHandler)

			// События
			admin.POST("/events", h.createEventHandler)
			admin.PUT("/events/:id", h.updateEventHandler)
			admin.DELETE("/events/:id", h.deleteEventHandler)

			// Учебные планы
			admin.POST("/plans", h.createPlanHandler)
			admin.PUT("/plans/:id", h.updatePlanHandler)
			admin.DELETE("/plans/:id", h.deletePlanHandler)

			// Доска почёта
			admin.POST("/honor", h.createHonorMemberHandler)
			admin.PUT("/honor/:id", h.updateHonorMemberHandler)
			admin.DELETE("/honor/:id", h.deleteHonorMemberHandler)

			// Достижения
			admin.POST("/achievements", h.createAchievementHandler)
			admin.PUT("/achievements/:id", h.updateAchievementHandler)
			admin.DELETE("/achievements/:id", h.deleteAchievementHandler)

			// Основатель
			admin.PUT("/founder", h.upsertFounderHandler)

			// Пользователи и роли
			admin.GET("/users", roleMiddleware("admin", "founder"), h.listUsersHandler)
			admin.PUT("/users/:id/role", roleMiddleware("admin", "founder"), h.updateUserRoleHandler)

			// Профили инструкторов
			admin.PUT("/instructors", roleMiddleware("instructor", "admin", "founder"), h.upsertInstructorProfileHandler)

			// База знаний
			admin.POST("/knowledge", roleMiddleware("instructor", "admin", "founder"), h.createKnowledgeHandler)
			admin.PUT("/knowledge/:id", roleMiddleware("instructor", "admin", "founder"), h.updateKnowledgeHandler)
			admin.DELETE("/knowledge/:id", roleMiddleware("instructor", "admin", "founder"), h.deleteKnowledgeHandler)

			// Глоссарий
			admin.POST("/glossary", roleMiddleware("instructor", "admin", "founder"), h.createGlossaryHandler)
			admin.PUT("/glossary/:id", roleMiddleware("instructor", "admin", "founder"), h.updateGlossaryHandler)
			admin.DELETE("/glossary/:id", roleMiddleware("instructor", "admin", "founder"), h.deleteGlossaryHandler)

			// Учебный прогресс
			admin.PUT("/progress", roleMiddleware("instructor", "admin", "founder"), h.upsertProgressHandler)
		}
	}

	// 404 для всех остальных маршрутов
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "маршрут не найден"})
	})

	return r
}
