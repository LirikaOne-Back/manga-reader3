// Исправленный интерфейс в файле internal/handler/manga.go
package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/gin-gonic/gin"
)

// MangaHandler обрабатывает HTTP-запросы, связанные с мангой
type MangaHandler struct {
	mangaService MangaService
	logger       *slog.Logger
}

// MangaService интерфейс сервиса манги
type MangaService interface {
	GetAll(ctx context.Context, filter domain.MangaFilter) ([]domain.Manga, int, error)
	GetByID(ctx context.Context, id int) (domain.Manga, error)
	Create(ctx context.Context, manga domain.Manga) (int, error)
	Update(ctx context.Context, manga domain.Manga) error
	Delete(ctx context.Context, id int) error
	GetGenres(ctx context.Context) ([]domain.Genre, error)
}

// NewMangaHandler создает новый экземпляр MangaHandler
func NewMangaHandler(mangaService MangaService, logger *slog.Logger) *MangaHandler {
	return &MangaHandler{
		mangaService: mangaService,
		logger:       logger,
	}
}

// Register регистрирует обработчики путей для манги
func (h *MangaHandler) Register(router *gin.RouterGroup) {
	manga := router.Group("/manga")
	{
		manga.GET("", h.getAllManga)
		manga.GET("/:id", h.getMangaByID)
		manga.POST("", h.createManga)
		manga.PUT("/:id", h.updateManga)
		manga.DELETE("/:id", h.deleteManga)
		manga.GET("/genres", h.getGenres)
	}
}

// getAllManga возвращает список манги с фильтрацией
// @Summary Получить список манги
// @Description Возвращает список манги с возможностью фильтрации и пагинации
// @Tags manga
// @Accept json
// @Produce json
// @Param genre query string false "Фильтр по жанру"
// @Param status query string false "Фильтр по статусу"
// @Param sortBy query string false "Поле для сортировки (title, rating, date)"
// @Param sortDesc query boolean false "Сортировка по убыванию"
// @Param search query string false "Поиск по названию"
// @Param page query int false "Номер страницы (начиная с 1)"
// @Param pageSize query int false "Размер страницы"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/manga [get]
func (h *MangaHandler) getAllManga(c *gin.Context) {
	h.logger.Info("handling get all manga request")

	var filter domain.MangaFilter

	// Получаем параметры фильтрации
	filter.Genre = c.Query("genre")
	filter.Status = c.Query("status")
	filter.SortBy = c.Query("sortBy")
	filter.SortDesc = c.Query("sortDesc") == "true"
	filter.Search = c.Query("search")

	// Параметры пагинации
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	filter.Page = page

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	filter.PageSize = pageSize

	// Получаем данные от сервиса
	mangas, total, err := h.mangaService.GetAll(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("failed to get manga list", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get manga list"})
		return
	}

	// Формируем ответ
	c.JSON(http.StatusOK, gin.H{
		"data":  mangas,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// getMangaByID возвращает мангу по идентификатору
// @Summary Получить мангу по ID
// @Description Возвращает детальную информацию о манге по её ID
// @Tags manga
// @Accept json
// @Produce json
// @Param id path int true "ID манги"
// @Success 200 {object} domain.Manga
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/manga/{id} [get]
func (h *MangaHandler) getMangaByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid manga id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid manga ID format"})
		return
	}

	manga, err := h.mangaService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get manga", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Manga not found"})
		return
	}

	c.JSON(http.StatusOK, manga)
}

// createManga создает новую мангу
// @Summary Создать новую мангу
// @Description Создает новую мангу с указанными данными
// @Tags manga
// @Accept json
// @Produce json
// @Param manga body domain.Manga true "Данные для создания манги"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/manga [post]
func (h *MangaHandler) createManga(c *gin.Context) {
	var manga domain.Manga

	if err := c.ShouldBindJSON(&manga); err != nil {
		h.logger.Error("invalid manga data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid manga data: " + err.Error()})
		return
	}

	id, err := h.mangaService.Create(c.Request.Context(), manga)
	if err != nil {
		h.logger.Error("failed to create manga", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to create manga: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Manga created successfully"})
}

// updateManga обновляет существующую мангу
// @Summary Обновить мангу
// @Description Обновляет существующую мангу по её ID
// @Tags manga
// @Accept json
// @Produce json
// @Param id path int true "ID манги"
// @Param manga body domain.Manga true "Данные для обновления манги"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/manga/{id} [put]
func (h *MangaHandler) updateManga(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid manga id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid manga ID format"})
		return
	}

	var manga domain.Manga
	if err := c.ShouldBindJSON(&manga); err != nil {
		h.logger.Error("invalid manga data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid manga data: " + err.Error()})
		return
	}

	// Устанавливаем ID из URL
	manga.ID = id

	err = h.mangaService.Update(c.Request.Context(), manga)
	if err != nil {
		h.logger.Error("failed to update manga", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to update manga: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manga updated successfully"})
}

// deleteManga удаляет мангу по идентификатору
// @Summary Удалить мангу
// @Description Удаляет мангу по её ID
// @Tags manga
// @Accept json
// @Produce json
// @Param id path int true "ID манги"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/manga/{id} [delete]
func (h *MangaHandler) deleteManga(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid manga id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid manga ID format"})
		return
	}

	err = h.mangaService.Delete(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to delete manga", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to delete manga: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manga deleted successfully"})
}

// getGenres возвращает список всех жанров
// @Summary Получить список жанров
// @Description Возвращает список всех доступных жанров манги
// @Tags manga
// @Accept json
// @Produce json
// @Success 200 {array} domain.Genre
// @Failure 500 {object} ErrorResponse
// @Router /api/manga/genres [get]
func (h *MangaHandler) getGenres(c *gin.Context) {
	genres, err := h.mangaService.GetGenres(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get genres", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get genres"})
		return
	}

	c.JSON(http.StatusOK, genres)
}
