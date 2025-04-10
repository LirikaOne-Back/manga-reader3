package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LirikaOne-Back/manga-reader3/internal/config"
	"github.com/LirikaOne-Back/manga-reader3/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализируем логгер
	log := logger.New(cfg.Logger)
	log.Info("Starting manga reader API server")

	// Устанавливаем режим Gin в зависимости от уровня логирования
	if cfg.Logger.Level == logger.LevelDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Инициализируем роутер
	router := gin.New()

	// Добавляем middleware для логирования
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware(log))

	// Настраиваем CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // В продакшене лучше указать конкретные домены
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Инициализируем API группу
	api := router.Group("/api")

	// Добавляем обработчики для различных ресурсов
	// TODO: Инициализировать репозитории, сервисы и обработчики
	// h.Register(api)

	// Простой эндпоинт для проверки работоспособности
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Создаем и настраиваем HTTP-сервер
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("Server started", "port", cfg.Server.Port)

	// Настраиваем корректное завершение приложения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Устанавливаем таймаут для завершения
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Пытаемся корректно завершить сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}

// loggerMiddleware создает middleware для логирования HTTP-запросов
func loggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Засекаем время начала запроса
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Обрабатываем запрос
		c.Next()

		// Логируем информацию о запросе
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// Определяем уровень логирования в зависимости от статус-кода
		if statusCode >= 500 {
			log.Error("Request processed",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", latency,
				"ip", clientIP,
				"error", c.Errors.String())
		} else if statusCode >= 400 {
			log.Warn("Request processed",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", latency,
				"ip", clientIP)
		} else {
			log.Info("Request processed",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", latency,
				"ip", clientIP)
		}
	}
}
