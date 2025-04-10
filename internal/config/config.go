package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/LirikaOne-Back/manga-reader3/pkg/logger"
	"github.com/joho/godotenv"
)

// Config содержит все настройки приложения
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Logger   logger.Config
	JWT      JWTConfig
	Storage  StorageConfig
	Redis    RedisConfig
}

// ServerConfig настройки HTTP-сервера
type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

// PostgresConfig настройки подключения к PostgreSQL
type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig настройки JWT-токенов
type JWTConfig struct {
	Secret           string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	SigningAlgorithm string
}

// StorageConfig настройки хранилища файлов
type StorageConfig struct {
	ImagesPath string
}

// RedisConfig содержит настройки подключения к Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
	TTL      time.Duration
}

// DSN возвращает строку подключения к PostgreSQL
func (pc PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pc.Username, pc.Password, pc.Host, pc.Port, pc.DBName, pc.SSLMode)
}

// Addr возвращает адрес подключения к Redis
func (rc RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", rc.Host, rc.Port)
}

// LoadConfig загружает конфигурацию из .env файла и переменных окружения
func LoadConfig(path string) (*Config, error) {
	// Загружаем .env файл, если он существует
	if err := godotenv.Load(path); err != nil {
		// Не считаем ошибку, если файла нет - используем переменные окружения
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Настройки сервера
	serverPort := getEnv("SERVER_PORT", "8080")
	shutdownTimeout, _ := strconv.Atoi(getEnv("SERVER_SHUTDOWN_TIMEOUT", "5"))
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))

	// Настройки PostgreSQL
	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPort := getEnv("POSTGRES_PORT", "5432")
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPass := getEnv("POSTGRES_PASSWORD", "postgres")
	pgDB := getEnv("POSTGRES_DB", "manga_reader")
	pgSSLMode := getEnv("POSTGRES_SSLMODE", "disable")

	// Настройки JWT
	jwtSecret := getEnv("JWT_SECRET", "super_secret_key")
	jwtAccessTTL, _ := strconv.Atoi(getEnv("JWT_ACCESS_TTL", "15"))  // минуты
	jwtRefreshTTL, _ := strconv.Atoi(getEnv("JWT_REFRESH_TTL", "7")) // дни
	jwtAlgorithm := getEnv("JWT_SIGNING_ALGORITHM", "HS256")

	// Настройки логгера
	logLevel := getEnv("LOG_LEVEL", "info")
	logJSON := getEnv("LOG_JSON", "false") == "true"

	// Настройки хранилища
	imagesPath := getEnv("STORAGE_IMAGES_PATH", "./data/images")

	// Настройки Redis
	redisHost := getEnv("REDIS_HOST", "redis")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	redisPoolSize, _ := strconv.Atoi(getEnv("REDIS_POOL_SIZE", "10"))
	redisTTL, _ := strconv.Atoi(getEnv("REDIS_TTL", "60")) // в минутах

	// Создаем и возвращаем конфигурацию
	return &Config{
		Server: ServerConfig{
			Port:            serverPort,
			ShutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
			ReadTimeout:     time.Duration(readTimeout) * time.Second,
			WriteTimeout:    time.Duration(writeTimeout) * time.Second,
		},
		Postgres: PostgresConfig{
			Host:     pgHost,
			Port:     pgPort,
			Username: pgUser,
			Password: pgPass,
			DBName:   pgDB,
			SSLMode:  pgSSLMode,
		},
		Logger: logger.Config{
			Level:      logger.Level(logLevel),
			JSONFormat: logJSON,
			Output:     os.Stdout,
			WithSource: true,
		},
		JWT: JWTConfig{
			Secret:           jwtSecret,
			AccessTokenTTL:   time.Duration(jwtAccessTTL) * time.Minute,
			RefreshTokenTTL:  time.Duration(jwtRefreshTTL) * 24 * time.Hour,
			SigningAlgorithm: jwtAlgorithm,
		},
		Storage: StorageConfig{
			ImagesPath: imagesPath,
		},
		Redis: RedisConfig{
			Host:     redisHost,
			Port:     redisPort,
			Password: redisPassword,
			DB:       redisDB,
			PoolSize: redisPoolSize,
			TTL:      time.Duration(redisTTL) * time.Minute,
		},
	}, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
