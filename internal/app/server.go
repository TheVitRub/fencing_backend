package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// shutdownTimeout — время ожидания graceful shutdown.
// Держим короче, чем terminationGracePeriodSeconds в k8s/docker.
const shutdownTimeout = 10 * time.Second

// Run запускает HTTP-сервер и блокируется до получения сигнала завершения
// (SIGINT / SIGTERM) или критической ошибки сервера.
// При остановке даёт активным запросам завершиться в течение shutdownTimeout.
func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%d", a.Config.Server.Host, a.Config.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		a.Logger.Info(a.Config.ServiceName+" HTTP-сервер запущен", a.Logger.F("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- fmt.Errorf("ошибка HTTP-сервера: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	var runErr error
	select {
	case <-quit:
		a.Logger.Info("получен сигнал завершения")
	case err := <-serverErrCh:
		runErr = err
		a.Logger.Error("критическая ошибка сервера", a.Logger.F("error", err))
	}

	// Graceful shutdown: новые соединения не принимаем, ждём завершения активных
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.Logger.Error("ошибка graceful shutdown HTTP-сервера", a.Logger.F("error", err))
	} else {
		a.Logger.Info("HTTP-сервер остановлен (graceful)")
	}

	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			a.Logger.Error("ошибка закрытия соединения с PostgreSQL", a.Logger.F("error", err))
		}
	}

	return runErr
}
