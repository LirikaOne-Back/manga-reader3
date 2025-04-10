package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LirikaOne-Back/manga-reader3/internal/config"
	"github.com/LirikaOne-Back/manga-reader3/internal/handler"
	"github.com/LirikaOne-Back/manga-reader3/internal/repository"
	"github.com/LirikaOne-Back/manga-reader3/internal/repository/postgres"
	"github.com/LirikaOne-Back/manga-reader3/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
	"github.com/redis/go-redis/v9"
)

// App представляет экземпляр приложения
type App struct {
	httpServer  *http.Server
	cfg         *config.Config
	logger      *slog.Logger
	db          *sqlx.DB
	redisClient *redis.Client
	handlers    *handler.Handler
}

// NewApp создает новый экземпляр приложения
func NewApp(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// Инициализируем подключение к PostgreSQL
	db, err := initPostgres(cfg.Postgres, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres: %w", err)
	}

	// Инициализируем подключение к Redis
	redisClient, err := initRedis(cfg.Redis, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	// Инициализируем репозитории
	repos := initRepositories(db, logger)

	// Инициализируем сервисы
	services := initServices(repos, cfg, redisClient, logger)

	// Инициализируем обработчики
	handlers := initHandlers(services, logger)

	// Инициализируем роутер
	router := initRouter(handlers, cfg, logger)

	// Инициализируем HTTP-сервер
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &App{
		httpServer:  httpServer,
		cfg:         cfg,
		logger:      logger,
		db:          db,
		redisClient: redisClient,
		handlers:    handlers,
	}, nil
}

// Run запускает приложение
func (a *App) Run() error {
	// Запускаем HTTP-сервер в отдельной горутине
	go func() {
		a.logger.Info("Starting HTTP server", "port", a.cfg.Server.Port)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("Failed to start HTTP server", "error", err)
		}
	}()

	// Создаем канал для сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнал завершения
	<-quit
	a.logger.Info("Shutting down server...")

	// Создаем контекст с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
	defer cancel()

	// Закрываем HTTP-сервер
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("Failed to shutdown HTTP server", "error", err)
		return err
	}

	// Закрываем соединение с PostgreSQL
	if err := a.db.Close(); err != nil {
		a.logger.Error("Failed to close PostgreSQL connection", "error", err)
		return err
	}

	// Закрываем соединение с Redis
	if err := a.redisClient.Close(); err != nil {
		a.logger.Error("Failed to close Redis connection", "error", err)
		return err
	}

	a.logger.Info("Server gracefully stopped")
	return nil
}

// initPostgres инициализирует подключение к PostgreSQL
func initPostgres(cfg config.PostgresConfig, logger *slog.Logger) (*sqlx.DB, error) {
	logger.Info("Connecting to PostgreSQL", "host", cfg.Host, "port", cfg.Port, "dbname", cfg.DBName)

	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	logger.Info("Connected to PostgreSQL successfully")
	return db, nil
}

// initRedis инициализирует подключение к Redis
func initRedis(cfg config.RedisConfig, logger *slog.Logger) (*redis.Client, error) {
	logger.Info("Connecting to Redis", "host", cfg.Host, "port", cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 10,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	})

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	logger.Info("Connected to Redis successfully")
	return client, nil
}

// initRepositories инициализирует репозитории
func initRepositories(db *sqlx.DB, logger *slog.Logger) *repository.Repositories {
	return &repository.Repositories{
		Manga:   postgres.NewMangaRepo(db, logger),
		Chapter: postgres.NewChapterRepo(db, logger),
		User:    postgres.NewUserRepo(db, logger),
	}
}

// initServices инициализирует сервисы
func initServices(
	repos *repository.Repositories,
	cfg *config.Config,
	redisClient *redis.Client,
	logger *slog.Logger,
) *service.Services {
	mangaService := service.NewMangaService(repos.Manga, logger)

	chapterService := service.NewChapterService(
		repos.Chapter,
		repos.Manga,
		logger,
		cfg.Storage.ImagesPath,
	)

	authConfig := service.AuthConfig{
		JWTSecret:    cfg.JWT.Secret,
		AccessTTL:    cfg.JWT.AccessTokenTTL,
		RefreshTTL:   cfg.JWT.RefreshTokenTTL,
		JWTAlgorithm: cfg.JWT.SigningAlgorithm,
		RedisClient:  redisClient,
	}

	authService := service.NewAuthService(repos.User, logger, authConfig)

	userService := service.NewUserService(repos.User, logger)

	return &service.Services{
		Manga:   mangaService,
		Chapter: chapterService,
		Auth:    authService,
		User:    userService,
	}
}

// initHandlers инициализирует обработчики
func initHandlers(
	services *service.Services,
	logger *slog.Logger,
) *handler.Handler {
	return handler.NewHandler(services, logger)
}

// initRouter инициализирует роутер Gin
func initRouter(
	h *handler.Handler,
	cfg *config.Config,
	logger *slog.Logger,
) *gin.Engine {
	// Настраиваем режим Gin в зависимости от уровня логирования
	if cfg.Logger.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Создаем роутер
	router := gin.New()

	// Регистрируем общие middleware
	router.Use(gin.Recovery())

	// Регистрируем обработчики
	h.Register(router)

	return router
}
