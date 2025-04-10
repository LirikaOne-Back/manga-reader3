package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/LirikaOne-Back/manga-reader3/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler объединяет все обработчики HTTP
type Handler struct {
	services   *service.Services
	logger     *slog.Logger
	manga      *MangaHandler
	chapter    *ChapterHandler
	auth       *AuthHandler
	user       *UserHandler
	middleware *Middleware
}

// NewHandler создает новый экземпляр Handler
func NewHandler(services *service.Services, logger *slog.Logger) *Handler {
	// Инициализируем middleware
	middleware := NewMiddleware(services.Auth, logger)

	// Инициализируем обработчики
	mangaHandler := NewMangaHandler(services.Manga, logger)
	chapterHandler := NewChapterHandler(services.Chapter, logger)
	authHandler := NewAuthHandler(services.Auth, logger)
	userHandler := NewUserHandler(services.User, middleware, logger)

	return &Handler{
		services:   services,
		logger:     logger,
		manga:      mangaHandler,
		chapter:    chapterHandler,
		auth:       authHandler,
		user:       userHandler,
		middleware: middleware,
	}
}

// Register регистрирует все обработчики HTTP
func (h *Handler) Register(router *gin.Engine) {
	// Добавляем middleware
	router.Use(h.middleware.Logger())
	router.Use(h.middleware.Recover())
	router.Use(h.middleware.CORS())
	router.Use(h.middleware.ContentTypeJSON())

	// Статические файлы для изображений
	router.Static("/images", "./data/images")

	// Простой эндпоинт для проверки работоспособности
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Группа API
	api := router.Group("/api")
	{
		// Регистрируем обработчики
		h.manga.Register(api)
		h.chapter.Register(api)
		h.auth.Register(api)
		h.user.Register(api)
	}

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
