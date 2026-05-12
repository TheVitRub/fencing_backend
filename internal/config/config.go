package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config — полная конфигурация сервиса, загружаемая из YAML + переменных окружения.
type Config struct {
	Environment string         // development | production
	ServiceName string         // fencing-club
	Server      ServerConfig   // HTTP-сервер
	Database    DatabaseConfig // PostgreSQL
	Auth        AuthConfig     // JWT
	Logging     LoggingConfig  // структурированное логирование
}

// ServerConfig содержит параметры HTTP-сервера.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig содержит параметры подключения к PostgreSQL.
type DatabaseConfig struct {
	Host                          string        `mapstructure:"host"`
	Port                          int           `mapstructure:"port"`
	User                          string        `mapstructure:"user"`
	Password                      string        `mapstructure:"password"`
	Name                          string        `mapstructure:"name"`
	SSLMode                       string        `mapstructure:"ssl_mode"`
	MaxOpenConns                  int           `mapstructure:"max_open_conns"`
	MaxIdleConns                  int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime               time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime               time.Duration `mapstructure:"conn_max_idle_time"`
	ConnectTimeout                time.Duration `mapstructure:"connect_timeout"`
	MaxRetries                    int           `mapstructure:"max_retries"`
	RetryInterval                 time.Duration `mapstructure:"retry_interval"`
	HealthCheckPeriod             time.Duration `mapstructure:"health_check_period"`
	HealthCheckReconnectThreshold int           `mapstructure:"health_check_reconnect_threshold"`
}

// AuthConfig содержит параметры JWT-аутентификации.
type AuthConfig struct {
	// JWTSecret — секретный ключ подписи токенов.
	// Передаётся через переменную окружения JWT_SECRET.
	JWTSecret string        `mapstructure:"jwt_secret"`
	JWTTTL    time.Duration `mapstructure:"jwt_ttl"`
}

// LoggingConfig содержит параметры логирования.
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load читает конфигурацию из файла configs/<env>.yaml и переменных окружения.
// Переменные окружения всегда перекрывают значения файла (удобно для секретов в k8s).
func Load() (*Config, error) {
	env := getEnvOrDefault("ENVIRONMENT", "development")
	serviceName := getEnvOrDefault("SERVICE_NAME", "fencing-club")

	v := viper.New()
	configFile := configFileName(env)
	v.SetConfigName(configFile)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath("../../configs")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("ошибка чтения конфигурации %s: %w", configFile, err)
		}
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnvSecrets(v)

	cfg := &Config{}
	cfg.Environment = getOrDefault(v.GetString("environment"), env)
	cfg.ServiceName = getOrDefault(v.GetString("service_name"), serviceName)

	// server
	cfg.Server.Host = getOrDefault(v.GetString("server.host"), "0.0.0.0")
	cfg.Server.Port = v.GetInt("server.port")

	// database
	cfg.Database.Host = v.GetString("database.host")
	cfg.Database.Port = v.GetInt("database.port")
	cfg.Database.User = v.GetString("database.user")
	cfg.Database.Password = v.GetString("database.password")
	cfg.Database.Name = v.GetString("database.name")
	cfg.Database.SSLMode = getOrDefault(v.GetString("database.ssl_mode"), "disable")
	cfg.Database.MaxOpenConns = getIntOrDefault(v.GetInt("database.max_open_conns"), 25)
	cfg.Database.MaxIdleConns = getIntOrDefault(v.GetInt("database.max_idle_conns"), 10)
	cfg.Database.ConnMaxLifetime = getDurationOrDefault(v.GetDuration("database.conn_max_lifetime"), 30*time.Minute)
	cfg.Database.ConnMaxIdleTime = getDurationOrDefault(v.GetDuration("database.conn_max_idle_time"), 10*time.Minute)
	cfg.Database.ConnectTimeout = getDurationOrDefault(v.GetDuration("database.connect_timeout"), 10*time.Second)
	cfg.Database.MaxRetries = getIntOrDefault(v.GetInt("database.max_retries"), 3)
	cfg.Database.RetryInterval = getDurationOrDefault(v.GetDuration("database.retry_interval"), 2*time.Second)
	cfg.Database.HealthCheckPeriod = getDurationOrDefault(v.GetDuration("database.health_check_period"), 30*time.Second)
	cfg.Database.HealthCheckReconnectThreshold = getIntOrDefault(v.GetInt("database.health_check_reconnect_threshold"), 3)

	// auth
	cfg.Auth.JWTSecret = v.GetString("auth.jwt_secret")
	cfg.Auth.JWTTTL = getDurationOrDefault(v.GetDuration("auth.jwt_ttl"), 24*time.Hour)

	// logging
	cfg.Logging.Level = getOrDefault(v.GetString("logging.level"), "info")
	cfg.Logging.Format = getOrDefault(v.GetString("logging.format"), "json")

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return cfg, nil
}

// bindEnvSecrets биндит переменные окружения для секретов.
// Это позволяет передавать чувствительные данные через Kubernetes Secrets,
// не записывая их в конфигурационные файлы.
func bindEnvSecrets(v *viper.Viper) {
	// server
	_ = v.BindEnv("server.host", "SERVER_HOST")
	_ = v.BindEnv("server.port", "SERVER_PORT", "PORT")

	// database secrets
	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.user", "DB_USER")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.name", "DB_NAME")
	_ = v.BindEnv("database.ssl_mode", "DB_SSL_MODE")

	// auth secrets
	_ = v.BindEnv("auth.jwt_secret", "JWT_SECRET")
	_ = v.BindEnv("auth.jwt_ttl", "JWT_TTL")

	// logging
	_ = v.BindEnv("logging.level", "LOGGING_LEVEL")
	_ = v.BindEnv("logging.format", "LOGGING_FORMAT")
}

// validate проверяет, что все обязательные поля конфигурации заполнены.
func (c *Config) validate() error {
	if c.Server.Port <= 0 {
		return fmt.Errorf("требуется переменная SERVER_PORT")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("требуется переменная DB_HOST")
	}
	if c.Database.Port <= 0 {
		return fmt.Errorf("требуется переменная DB_PORT")
	}
	if c.Database.User == "" {
		return fmt.Errorf("требуется переменная DB_USER")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("требуется переменная DB_PASSWORD")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("требуется переменная DB_NAME")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("требуется переменная JWT_SECRET")
	}

	// В production запрещаем небезопасный SSL-режим
	if c.IsProduction() && c.Database.SSLMode == "disable" {
		return fmt.Errorf("в окружении production нельзя использовать database.ssl_mode='disable'")
	}

	return nil
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// IsDevelopment возвращает true, если сервис запущен в режиме разработки.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction возвращает true, если сервис запущен в боевом окружении.
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// configFileName возвращает имя файла конфигурации для заданного окружения.
func configFileName(env string) string {
	switch env {
	case "production", "prod":
		return "prod"
	default:
		return "dev"
	}
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getOrDefault(v, def string) string {
	if v != "" {
		return v
	}
	return def
}

func getIntOrDefault(v, def int) int {
	if v > 0 {
		return v
	}
	return def
}

func getDurationOrDefault(v, def time.Duration) time.Duration {
	if v > 0 {
		return v
	}
	return def
}
