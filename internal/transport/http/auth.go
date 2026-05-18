package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"fencing-club/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// loginHandler обрабатывает POST /api/auth/login.
// Принимает логин/пароль, возвращает JWT-токен при успехе.
func (h *Handler) loginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) registerHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса: " + err.Error()})
		return
	}
	resp, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) meHandler(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		return
	}
	resp, err := h.svc.Me(c.Request.Context(), claims.UserID)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) oauthStartHandler(c *gin.Context) {
	provider := c.Param("provider")
	redirectURI := oauthRedirectURI(c, provider)
	target, err := h.svc.OAuthStartURL(provider, redirectURI)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.Redirect(http.StatusFound, target)
}

func (h *Handler) oauthCallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	resp, err := h.svc.OAuthCallback(c.Request.Context(), provider, c.Query("code"), oauthRedirectURI(c, provider))
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	tokenJSON, _ := json.Marshal(resp.Token)
	userJSON, _ := json.Marshal(resp.User)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!doctype html><meta charset="utf-8"><script>
localStorage.setItem('fc_token', %s);
localStorage.setItem('fc_user', JSON.stringify(%s));
location.href = '/';
</script>`, string(tokenJSON), string(userJSON))
}

func oauthRedirectURI(c *gin.Context, provider string) string {
	if base := os.Getenv("PUBLIC_BASE_URL"); base != "" {
		u, err := url.Parse(base)
		if err == nil {
			u.Path = "/api/auth/oauth/" + provider + "/callback"
			u.RawQuery = ""
			return u.String()
		}
	}
	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	return scheme + "://" + c.Request.Host + "/api/auth/oauth/" + provider + "/callback"
}
