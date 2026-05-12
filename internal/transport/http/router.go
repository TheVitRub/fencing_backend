package http

import (
	"fencing-club/internal/application"
	"fencing-club/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Handler — корневой объект HTTP-слоя.
// Содержит ссылку на сервис приложения и логгер.
type Handler struct {
	svc    *application.Service
	logger *logger.Logger
}

// NewRouter создаёт Gin-роутер с зарегистрированными маршрутами.
// Публичные маршруты доступны без токена; маршруты /admin требуют JWT.
func NewRouter(svc *application.Service, log *logger.Logger) *gin.Engine {
	h := &Handler{svc: svc, logger: log}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		// --- Аутентификация ---
		auth := api.Group("/auth")
		auth.POST("/login", h.loginHandler)

		// --- Публичные маршруты ---
		api.GET("/pages/:slug", h.getPageHandler)
		api.GET("/events", h.listEventsHandler)
		api.GET("/plans", h.listPlansHandler)
		api.GET("/honor", h.listHonorMembersHandler)
		api.GET("/achievements", h.listAchievementsHandler)
		api.GET("/founder", h.getFounderHandler)

		// --- Административные маршруты (требуют JWT) ---
		admin := api.Group("/admin")
		admin.Use(authMiddleware(svc))
		{
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
		}
	}

	return r
}
