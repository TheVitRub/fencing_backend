// Package pgerror содержит утилиты для работы с ошибками PostgreSQL.
package pgerror

import (
	"errors"

	"github.com/lib/pq"
)

// Коды ошибок PostgreSQL, по которым стоит делать retry.
const (
	CodeDeadlock              = "40P01" // deadlock detected
	CodeSerializationFailure  = "40001" // could not serialize access
	CodeConnectionException   = "08000" // connection exception
	CodeConnectionDoesntExist = "08003" // connection does not exist
	CodeConnectionFailure     = "08006" // connection failure
)

// IsRetryable возвращает true, если ошибка является транзитной и операцию можно повторить.
func IsRetryable(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case CodeDeadlock, CodeSerializationFailure,
			CodeConnectionException, CodeConnectionDoesntExist, CodeConnectionFailure:
			return true
		}
	}
	return false
}

// GetErrorName возвращает имя ошибки PostgreSQL для логирования.
func GetErrorName(err error) string {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return string(pgErr.Code)
	}
	return "unknown"
}
