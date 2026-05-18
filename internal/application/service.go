// Package application реализует бизнес-логику сервиса.
//
// Прикладной слой (application) знает только об интерфейсах репозиториев
// и доменных типах. Он не знает об HTTP, базах данных или конкретных драйверах.
//
// Каждая сущность реализована в своём файле для удобства навигации:
//   - auth.go       — авторизация администратора, JWT
//   - pages.go      — управление страницами сайта
//   - events.go     — события клуба
//   - plans.go      — учебные планы
//   - honor.go      — доска почёта
//   - achievements.go — достижения школы
//   - founder.go    — информация об основателе
package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"fencing-club/internal/domain/apperrors"
	"fencing-club/internal/domain/interfaces"
	"fencing-club/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
)

// Service — центральный объект прикладного слоя.
// Все операции выполняются через него, что упрощает подключение
// middleware (трассировка, кэш) в одном месте.
type Service struct {
	repo      interfaces.DBRepository
	jwtSecret string
	jwtTTL    time.Duration
	logger    *logger.Logger
}

// New создаёт Service с переданными зависимостями.
func New(
	repo interfaces.DBRepository,
	jwtSecret string,
	jwtTTL time.Duration,
	log *logger.Logger,
) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
		logger:    log,
	}
}

// --- JWT / внутренние утилиты -----------------------------------------

// Claims — полезная нагрузка JWT-токена администратора.
type Claims struct {
	AdminID int64  `json:"admin_id,omitempty"` // legacy alias
	UserID  int64  `json:"user_id"`
	Login   string `json:"login"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// issueToken выпускает JWT для администратора.
func (s *Service) issueToken(userID int64, login, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		AdminID: userID,
		UserID:  userID,
		Login:   login,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("не удалось подписать токен: %w", err)
	}
	return signed, nil
}

// ValidateToken проверяет JWT и возвращает Claims.
// При любой ошибке возвращает apperrors.KindUnauthorized.
func (s *Service) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.Unauthorizedf("неожиданный метод подписи токена")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, apperrors.Unauthorizedf("недействительный токен")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, apperrors.Unauthorizedf("не удалось разобрать claims токена")
	}
	return claims, nil
}

// wrapRepoError оборачивает ошибки репозитория в доменные ошибки.
// sql.ErrNoRows → apperrors.KindNotFound, всё остальное → apperrors.KindInternal.
func wrapRepoError(op string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return apperrors.NotFoundf("%s: запись не найдена", op)
	}
	return apperrors.New(apperrors.KindInternal, fmt.Sprintf("%s: внутренняя ошибка", op), err)
}

// withContext — тип-обёртка для передачи контекста в методы, которые его не принимают явно.
// Удобен для вызовов внутри горутин при добавлении трассировки в будущем.
type withContext struct {
	ctx context.Context //nolint:containedctx
}
