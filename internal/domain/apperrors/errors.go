// Package apperrors определяет типизированные доменные ошибки сервиса.
//
// Принцип работы: прикладной слой (application) возвращает *Error с нужным Kind,
// транспортный слой (transport/http) читает Kind и преобразует его в HTTP-статус.
// Это позволяет менять транспорт (HTTP → gRPC) без правок бизнес-логики.
package apperrors

import "fmt"

// Kind описывает категорию ошибки, не зависящую от транспорта.
type Kind string

const (
	// KindValidation — некорректные входные данные.
	KindValidation Kind = "validation"
	// KindNotFound — запрошенный ресурс не найден.
	KindNotFound Kind = "not_found"
	// KindConflict — нарушение уникальности / конкурентное изменение.
	KindConflict Kind = "conflict"
	// KindUnauthorized — не авторизован (неверные учётные данные).
	KindUnauthorized Kind = "unauthorized"
	// KindForbidden — доступ запрещён (нет прав).
	KindForbidden Kind = "forbidden"
	// KindInternal — внутренняя ошибка сервиса.
	KindInternal Kind = "internal"
)

// Error — доменная ошибка с категорией и опциональной причиной.
type Error struct {
	kind  Kind
	msg   string
	cause error
}

func (e *Error) Error() string { return e.msg }
func (e *Error) Unwrap() error { return e.cause }
func (e *Error) Kind() Kind    { return e.kind }

// New создаёт ошибку заданной категории с сообщением и причиной.
// cause может быть nil, если оборачивать нечего.
func New(kind Kind, msg string, cause error) error {
	return &Error{kind: kind, msg: msg, cause: cause}
}

// Validationf создаёт ошибку валидации с форматированным сообщением.
func Validationf(format string, args ...any) error {
	return New(KindValidation, fmt.Sprintf(format, args...), nil)
}

// NotFoundf создаёт ошибку «не найдено» с форматированным сообщением.
func NotFoundf(format string, args ...any) error {
	return New(KindNotFound, fmt.Sprintf(format, args...), nil)
}

// Conflictf создаёт ошибку конфликта с форматированным сообщением.
func Conflictf(format string, args ...any) error {
	return New(KindConflict, fmt.Sprintf(format, args...), nil)
}

// Unauthorizedf создаёт ошибку авторизации с форматированным сообщением.
func Unauthorizedf(format string, args ...any) error {
	return New(KindUnauthorized, fmt.Sprintf(format, args...), nil)
}

// Forbiddenf создаёт ошибку запрета доступа с форматированным сообщением.
func Forbiddenf(format string, args ...any) error {
	return New(KindForbidden, fmt.Sprintf(format, args...), nil)
}

// Internalf создаёт внутреннюю ошибку с форматированным сообщением.
func Internalf(format string, args ...any) error {
	return New(KindInternal, fmt.Sprintf(format, args...), nil)
}
