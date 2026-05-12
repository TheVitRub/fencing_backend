package app

import (
	"fmt"

	"fencing-club/internal/application"
	"fencing-club/internal/config"
	infraRepo "fencing-club/internal/infrastructure/repository"
	dbRepository "fencing-club/internal/repository/db"
	httpTransport "fencing-club/internal/transport/http"
	"fencing-club/pkg/logger"

	"github.com/gin-gonic/gin"
)

// App — корневой объект приложения.
// Хранит все зависимости и предоставляет метод Run для запуска сервера.
type App struct {
	Config  *config.Config
	Logger  *logger.Logger
	DB      *infraRepo.PostgresDB
	Service *application.Service
	Router  *gin.Engine
}

// New инициализирует все зависимости приложения через DI:
// конфиг → логгер → postgres → repository → service → router.
func New(cfg *config.Config) (*App, error) {
	log := logger.Init(cfg.ServiceName, cfg.Environment)

	db, err := infraRepo.NewPostgres(cfg.Database, log)
	if err != nil {
		log.Error("ошибка создания postgres клиента", log.F("error", err))
		return nil, fmt.Errorf("ошибка создания postgres клиента: %w", err)
	}
	log.Info("подключение к PostgreSQL установлено")

	repo := dbRepository.New(db, log)
	svc := application.New(repo, cfg.Auth.JWTSecret, cfg.Auth.JWTTTL, log)
	router := httpTransport.NewRouter(svc, log)

	return &App{
		Config:  cfg,
		Logger:  log,
		DB:      db,
		Service: svc,
		Router:  router,
	}, nil
}
