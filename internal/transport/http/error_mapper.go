package http

import (
	"errors"

	"fencing-club/internal/domain/apperrors"
)

// mapError преобразует доменную ошибку в HTTP-статус и сообщение для клиента.
// Внутренние детали ошибок (apperrors.KindInternal) скрываются от клиента.
func mapError(err error) (status int, message string) {
	var appErr *apperrors.Error
	if errors.As(err, &appErr) {
		switch appErr.Kind() {
		case apperrors.KindValidation:
			return 400, appErr.Error()
		case apperrors.KindNotFound:
			return 404, appErr.Error()
		case apperrors.KindConflict:
			return 409, appErr.Error()
		case apperrors.KindUnauthorized:
			return 401, appErr.Error()
		case apperrors.KindForbidden:
			return 403, appErr.Error()
		}
	}
	return 500, "внутренняя ошибка сервиса"
}
