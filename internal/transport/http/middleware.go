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

func getClaims(c *gin.Context) (*application.Claims, bool) {
	raw, ok := c.Get("claims")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "нужен вход"})
		return nil, false
	}
	claims, ok := raw.(*application.Claims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "некорректный токен"})
		return nil, false
	}
	if claims.UserID == 0 {
		claims.UserID = claims.AdminID
	}
	return claims, true
}

func roleMiddleware(allowed ...string) gin.HandlerFunc {
	allowedSet := map[string]bool{}
	for _, role := range allowed {
		allowedSet[role] = true
	}
	return func(c *gin.Context) {
		claims, ok := getClaims(c)
		if !ok {
			return
		}
		if !allowedSet[claims.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "недостаточно прав"})
			return
		}
		c.Next()
	}
}
