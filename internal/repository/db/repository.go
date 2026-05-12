// Package db реализует интерфейс interfaces.DBRepository поверх PostgreSQL.
//
// DBRepository содержит только зависимости (пул соединений + логгер),
// а каждая сущность (admins, events, …) реализована в своём файле.
// Такое разбиение позволяет легко найти нужный метод и не раздувать один файл.
package db

import (
	infraRepo "fencing-club/internal/infrastructure/repository"
	"fencing-club/pkg/logger"
)

// DBRepository реализует interfaces.DBRepository.
type DBRepository struct {
	db  *infraRepo.PostgresDB
	log *logger.Logger
}

// New создаёт новый DBRepository.
func New(db *infraRepo.PostgresDB, log *logger.Logger) *DBRepository {
	return &DBRepository{db: db, log: log}
}
