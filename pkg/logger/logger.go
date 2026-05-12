// Package logger предоставляет обёртку над go.uber.org/zap с удобным API.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger — обёртка над *zap.Logger с методом F для построения полей.
type Logger struct {
	*zap.Logger
}

// Init создаёт новый логгер для указанного сервиса и окружения.
// В режиме development используется консольный вывод с цветами,
// в production — JSON-формат для сборщиков логов.
func Init(serviceName, environment string) *Logger {
	var cfg zap.Config

	if environment == "development" || environment == "dev" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	base, err := cfg.Build(
		zap.AddCallerSkip(0),
		zap.Fields(
			zap.String("service", serviceName),
			zap.String("env", environment),
		),
	)
	if err != nil {
		// Если логгер не удалось создать — падаем: без логов работать нельзя.
		panic("не удалось инициализировать логгер: " + err.Error())
	}

	return &Logger{base}
}

// F — сокращение для zap.Any, чтобы не импортировать zap во всех пакетах.
func (l *Logger) F(key string, val any) zap.Field {
	return zap.Any(key, val)
}
