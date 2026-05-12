package application

import (
	"context"

	"fencing-club/internal/domain/apperrors"
	"fencing-club/internal/domain/dto"

	"golang.org/x/crypto/bcrypt"
)

// Login проверяет учётные данные администратора и возвращает JWT-токен.
// Возвращает apperrors.KindUnauthorized при неверном логине или пароле.
func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	admin, err := s.repo.GetAdminByLogin(ctx, req.Login)
	if err != nil {
		// Скрываем детали: не сообщаем, существует ли логин.
		return nil, apperrors.Unauthorizedf("неверный логин или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.Unauthorizedf("неверный логин или пароль")
	}

	token, err := s.issueToken(admin.ID, admin.Login)
	if err != nil {
		return nil, apperrors.New(apperrors.KindInternal, "не удалось выпустить токен", err)
	}

	s.logger.Info("администратор вошёл в систему",
		s.logger.F("admin_id", admin.ID),
		s.logger.F("login", admin.Login),
	)
	return &dto.LoginResponse{Token: token}, nil
}
