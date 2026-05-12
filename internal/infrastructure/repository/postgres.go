// Package repository содержит инфраструктурные примитивы для работы с PostgreSQL.
//
// PostgresDB — обёртка над *sql.DB и *sqlx.DB, которая добавляет:
//   - повторные попытки (retry) при транзитных ошибках (deadlock, connection lost);
//   - фоновую проверку соединения (health check);
//   - единое место для настройки connection pool.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"fencing-club/internal/config"
	"fencing-club/pkg/logger"
	"fencing-club/pkg/pgerror"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // регистрируем драйвер postgres
)

// PostgresDB — потокобезопасная обёртка над пулом соединений PostgreSQL.
type PostgresDB struct {
	db  *sql.DB
	dbx *sqlx.DB

	cfg    config.DatabaseConfig
	logger *logger.Logger

	healthTicker *time.Ticker
	stopHealthCh chan struct{}
	closeOnce    sync.Once

	mu        sync.RWMutex
	isHealthy bool
	failCount int
}

// Executor — интерфейс для методов, работающих как с пулом, так и внутри транзакции.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// NewPostgres создаёт PostgresDB: подключается к БД с повторными попытками
// и запускает фоновый health-check.
func NewPostgres(cfg config.DatabaseConfig, log *logger.Logger) (*PostgresDB, error) {
	db, err := connectWithRetry(cfg, log)
	if err != nil {
		return nil, err
	}

	pg := &PostgresDB{
		db:           db,
		dbx:          sqlx.NewDb(db, "postgres"),
		cfg:          cfg,
		logger:       log,
		stopHealthCh: make(chan struct{}),
		isHealthy:    true,
	}
	pg.startHealthCheck()
	return pg, nil
}

// --- Подключение -------------------------------------------------------

func connectWithRetry(cfg config.DatabaseConfig, log *logger.Logger) (*sql.DB, error) {
	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 5
	}
	retryInterval := cfg.RetryInterval
	if retryInterval <= 0 {
		retryInterval = 2 * time.Second
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err := connectOnce(cfg)
		if err == nil {
			return db, nil
		}
		lastErr = err
		log.Warn("попытка подключения к PostgreSQL не удалась",
			log.F("attempt", attempt),
			log.F("max_retries", maxRetries),
			log.F("error", err),
		)
		time.Sleep(retryInterval)
	}
	return nil, fmt.Errorf("не удалось подключиться к PostgreSQL после %d попыток: %w", maxRetries, lastErr)
}

func connectOnce(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть соединение: %w", err)
	}

	timeout := cfg.ConnectTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping не прошёл: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	return db, nil
}

// --- Health-check ------------------------------------------------------

func (p *PostgresDB) startHealthCheck() {
	period := p.cfg.HealthCheckPeriod
	if period <= 0 {
		period = 30 * time.Second
	}
	p.healthTicker = time.NewTicker(period)

	go func() {
		for {
			select {
			case <-p.healthTicker.C:
				p.checkHealth()
			case <-p.stopHealthCh:
				p.healthTicker.Stop()
				return
			}
		}
	}()
}

func (p *PostgresDB) checkHealth() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	err := db.PingContext(ctx)

	p.mu.Lock()
	defer p.mu.Unlock()

	wasHealthy := p.isHealthy
	if err != nil {
		p.failCount++
		p.isHealthy = false
		if wasHealthy {
			p.logger.Error("PostgreSQL стал недоступен", p.logger.F("error", err))
		}
		return
	}

	p.failCount = 0
	p.isHealthy = true
	if !wasHealthy {
		p.logger.Info("PostgreSQL снова доступен",
			p.logger.F("open_conns", db.Stats().OpenConnections),
		)
	}
}

// IsHealthy возвращает текущее состояние соединения (по результату последнего health-check).
func (p *PostgresDB) IsHealthy() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isHealthy
}

// Ping выполняет немедленную проверку доступности БД.
func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.getDB().PingContext(ctx)
}

// Close закрывает соединение и останавливает health-check горутину.
func (p *PostgresDB) Close() error {
	p.closeOnce.Do(func() { close(p.stopHealthCh) })
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.db.Close()
}

// Stats возвращает статистику connection pool для health-endpoint'а.
func (p *PostgresDB) Stats() sql.DBStats {
	return p.getDB().Stats()
}

// --- Запросы с retry ---------------------------------------------------

// Exec выполняет INSERT / UPDATE / DELETE.
// При tx != nil — выполняет внутри транзакции (retry отключён).
func (p *PostgresDB) Exec(ctx context.Context, tx Executor, query string, args ...any) (sql.Result, error) {
	exec := p.executorOrDB(tx)
	if tx != nil {
		return exec.ExecContext(ctx, query, args...)
	}
	return p.retryExec(ctx, query, args...)
}

func (p *PostgresDB) retryExec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		res, err := p.getDB().ExecContext(ctx, query, args...)
		if err == nil {
			return res, nil
		}
		lastErr = err
		if !pgerror.IsRetryable(err) {
			return nil, err
		}
		p.logger.Warn("retry Exec",
			p.logger.F("attempt", attempt),
			p.logger.F("pg_code", pgerror.GetErrorName(err)),
		)
		if err := sleepCtx(ctx, time.Duration(attempt*100)*time.Millisecond); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("Exec: 3 попытки исчерпаны: %w", lastErr)
}

// QueryRow выполняет SELECT с одной строкой результата и сканирует в dest.
// При tx != nil — выполняет внутри транзакции.
func (p *PostgresDB) QueryRow(ctx context.Context, tx Executor, dest []any, query string, args ...any) error {
	exec := p.executorOrDB(tx)
	if tx != nil {
		return exec.QueryRowContext(ctx, query, args...).Scan(dest...)
	}
	return p.retryQueryRow(ctx, dest, query, args...)
}

func (p *PostgresDB) retryQueryRow(ctx context.Context, dest []any, query string, args ...any) error {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		err := p.getDB().QueryRowContext(ctx, query, args...).Scan(dest...)
		if err == nil {
			return nil
		}
		if errors.Is(err, sql.ErrNoRows) {
			return err // не ретраим — запись просто не существует
		}
		lastErr = err
		if !pgerror.IsRetryable(err) {
			return err
		}
		p.logger.Warn("retry QueryRow",
			p.logger.F("attempt", attempt),
			p.logger.F("pg_code", pgerror.GetErrorName(err)),
		)
		if err := sleepCtx(ctx, time.Duration(attempt*100)*time.Millisecond); err != nil {
			return err
		}
	}
	return fmt.Errorf("QueryRow: 3 попытки исчерпаны: %w", lastErr)
}

// Select выполняет SELECT и маппит результат в слайс структур через sqlx.
// При tx != nil — выполняет внутри транзакции.
func (p *PostgresDB) Select(ctx context.Context, tx Executor, dest any, query string, args ...any) error {
	if tx != nil {
		if etx, ok := tx.(sqlx.ExtContext); ok {
			return sqlx.SelectContext(ctx, etx, dest, query, args...)
		}
		// Стандартная *sql.Tx: оборачиваем в sqlx.Tx для корректного маппинга.
		if stdTx, ok := tx.(*sql.Tx); ok {
			p.mu.RLock()
			wrapped := sqlx.Tx{Tx: stdTx, Mapper: p.dbx.Mapper}
			p.mu.RUnlock()
			return sqlx.SelectContext(ctx, &wrapped, dest, query, args...)
		}
	}

	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		p.mu.RLock()
		dbx := p.dbx
		p.mu.RUnlock()

		err := dbx.SelectContext(ctx, dest, query, args...)
		if err == nil {
			return nil
		}
		lastErr = err
		if !pgerror.IsRetryable(err) {
			return err
		}
		if err := sleepCtx(ctx, time.Duration(attempt*100)*time.Millisecond); err != nil {
			return err
		}
	}
	return fmt.Errorf("Select: 3 попытки исчерпаны: %w", lastErr)
}

// Get выполняет SELECT и маппит одну строку результата в структуру через sqlx.
func (p *PostgresDB) Get(ctx context.Context, dest any, query string, args ...any) error {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		p.mu.RLock()
		dbx := p.dbx
		p.mu.RUnlock()

		err := dbx.GetContext(ctx, dest, query, args...)
		if err == nil {
			return nil
		}
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		lastErr = err
		if !pgerror.IsRetryable(err) {
			return err
		}
		if err := sleepCtx(ctx, time.Duration(attempt*100)*time.Millisecond); err != nil {
			return err
		}
	}
	return fmt.Errorf("Get: 3 попытки исчерпаны: %w", lastErr)
}

// --- helpers -----------------------------------------------------------

func (p *PostgresDB) getDB() *sql.DB {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.db
}

// executorOrDB возвращает tx если он передан, иначе — пул соединений.
func (p *PostgresDB) executorOrDB(tx Executor) Executor {
	if tx != nil {
		return tx
	}
	return p.getDB()
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
