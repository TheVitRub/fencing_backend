package http

import (
	"net/http"
	"strings"

	"fencing-club/internal/application"

	"github.com/gin-gonic/gin"
)

// corsMiddleware разрешает кросс-доменные запросы от фронтенда.
// В production рекомендуется ограничить AllowOrigins конкретным доменом.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// authMiddleware проверяет JWT-токен из заголовка Authorization: Bearer <token>.
// При успехе помещает claims в контекст gin под ключом "claims".
func authMiddleware(svc *application.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "отсутствует заголовок Authorization"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "неверный формат заголовка Authorization"})
			return
		}

		claims, err := svc.ValidateToken(parts[1])
		if err != nil {
			status, msg := mapError(err)
			c.AbortWithStatusJSON(status, gin.H{"error": msg})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
