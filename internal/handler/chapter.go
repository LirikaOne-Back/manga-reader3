package handler

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/gin-gonic/gin"
)

// ChapterHandler обрабатывает HTTP-запросы, связанные с главами манги
type ChapterHandler struct {
	chapterService ChapterService
	logger         *slog.Logger
}

// ChapterService интерфейс сервиса глав
type ChapterService interface {
	GetByMangaID(ctx gin.Context, mangaID int) ([]domain.Chapter, error)
	GetByID(ctx gin.Context, id int) (domain.Chapter, error)
	Create(ctx gin.Context, chapter domain.Chapter) (int, error)
	Update(ctx gin.Context, chapter domain.Chapter) error
	Delete(ctx gin.Context, id int) error
	GetPages(ctx gin.Context, chapterID int) ([]domain.Page, error)
	AddPage(ctx gin.Context, page domain.Page, imageData []byte) (int, error)
	DeletePage(ctx gin.Context, id int) error
}

// NewChapterHandler создает новый экземпляр ChapterHandler
func NewChapterHandler(chapterService ChapterService, logger *slog.Logger) *ChapterHandler {
	return &ChapterHandler{
		chapterService: chapterService,
		logger:         logger,
	}
}

// Register регистрирует обработчики путей для глав
func (h *ChapterHandler) Register(router *gin.RouterGroup) {
	chapters := router.Group("/chapters")
	{
		chapters.GET("/manga/:manga_id", h.getChaptersByManga)
		chapters.GET("/:id", h.getChapterByID)
		chapters.POST("", h.authMiddleware("moderator"), h.createChapter)
		chapters.PUT("/:id", h.authMiddleware("moderator"), h.updateChapter)
		chapters.DELETE("/:id", h.authMiddleware("moderator"), h.deleteChapter)

		// Пути для работы со страницами
		chapters.GET("/:id/pages", h.getChapterPages)
		chapters.POST("/:id/pages", h.authMiddleware("moderator"), h.addChapterPage)
		chapters.DELETE("/pages/:page_id", h.authMiddleware("moderator"), h.deleteChapterPage)
	}
}

// getChaptersByManga возвращает список глав для указанной манги
// @Summary Получить главы манги
// @Description Возвращает список глав для указанной манги
// @Tags chapters
// @Accept json
// @Produce json
// @Param manga_id path int true "ID манги"
// @Success 200 {array} domain.Chapter
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/chapters/manga/{manga_id} [get]
func (h *ChapterHandler) getChaptersByManga(c *gin.Context) {
	mangaID, err := strconv.Atoi(c.Param("manga_id"))
	if err != nil {
		h.logger.Error("invalid manga_id format", "manga_id", c.Param("manga_id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid manga ID format"})
		return
	}

	chapters, err := h.chapterService.GetByMangaID(*c, mangaID)
	if err != nil {
		h.logger.Error("failed to get chapters", "manga_id", mangaID, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get chapters: " + err.Error()})
		return
	}

	if len(chapters) == 0 {
		c.JSON(http.StatusOK, []domain.Chapter{})
		return
	}

	c.JSON(http.StatusOK, chapters)
}

// getChapterByID возвращает главу по ID
// @Summary Получить главу по ID
// @Description Возвращает детальную информацию о главе по её ID
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "ID главы"
// @Success 200 {object} domain.Chapter
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/chapters/{id} [get]
func (h *ChapterHandler) getChapterByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid chapter id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid chapter ID format"})
		return
	}

	chapter, err := h.chapterService.GetByID(*c, id)
	if err != nil {
		h.logger.Error("failed to get chapter", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Chapter not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, chapter)
}

// createChapter создает новую главу
// @Summary Создать новую главу
// @Description Создает новую главу с указанными данными
// @Tags chapters
// @Accept json
// @Produce json
// @Param chapter body domain.Chapter true "Данные для создания главы"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/chapters [post]
func (h *ChapterHandler) createChapter(c *gin.Context) {
	var chapter domain.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		h.logger.Error("invalid chapter data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid chapter data: " + err.Error()})
		return
	}

	id, err := h.chapterService.Create(*c, chapter)
	if err != nil {
		h.logger.Error("failed to create chapter", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to create chapter: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Chapter created successfully",
	})
}

// updateChapter обновляет существующую главу
// @Summary Обновить главу
// @Description Обновляет существующую главу по её ID
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "ID главы"
// @Param chapter body domain.Chapter true "Данные для обновления главы"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/chapters/{id} [put]
func (h *ChapterHandler) updateChapter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid chapter id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid chapter ID format"})
		return
	}

	var chapter domain.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		h.logger.Error("invalid chapter data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid chapter data: " + err.Error()})
		return
	}

	// Устанавливаем ID из URL
	chapter.ID = id

	err = h.chapterService.Update(*c, chapter)
	if err != nil {
		h.logger.Error("failed to update chapter", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to update chapter: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Chapter updated successfully",
	})
}

// deleteChapter удаляет главу по ID
// @Summary Удалить главу
// @Description Удаляет главу по её ID
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "ID главы"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/chapters/{id} [delete]
func (h *ChapterHandler) deleteChapter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid chapter id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid chapter ID format"})
		return
	}

	err = h.chapterService.Delete(*c, id)
	if err != nil {
		h.logger.Error("failed to delete chapter", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to delete chapter: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Chapter deleted successfully",
	})
}

// getChapterPages возвращает список страниц для указанной главы
// @Summary Получить страницы главы
// @Description Возвращает список страниц для указанной главы
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "ID главы"
// @Success 200 {array} domain.Page
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/chapters/{id}/pages [get]
func (h *ChapterHandler) getChapterPages(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid chapter id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid chapter ID format"})
		return
	}

	pages, err := h.chapterService.GetPages(*c, id)
	if err != nil {
		h.logger.Error("failed to get pages", "chapter_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get pages: " + err.Error()})
		return
	}

	if len(pages) == 0 {
		c.JSON(http.StatusOK, []domain.Page{})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// addChapterPage добавляет новую страницу в главу
// @Summary Добавить страницу
// @Description Добавляет новую страницу в главу
// @Tags chapters
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID главы"
// @Param number formData int false "Номер страницы (если не указан, будет добавлена в конец)"
// @Param image formData file true "Изображение страницы"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/chapters/{id}/pages [post]
func (h *ChapterHandler) addChapterPage(c *gin.Context) {
	chapterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid chapter id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid chapter ID format"})
		return
	}

	// Получаем номер страницы (необязательный параметр)
	numberStr := c.PostForm("number")
	var number int
	if numberStr != "" {
		number, err = strconv.Atoi(numberStr)
		if err != nil || number <= 0 {
			h.logger.Error("invalid page number", "number", numberStr)
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid page number"})
			return
		}
	}

	// Получаем изображение
	file, err := c.FormFile("image")
	if err != nil {
		h.logger.Error("failed to get image file", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Failed to get image file: " + err.Error()})
		return
	}

	// Ограничиваем размер файла (например, до 5 МБ)
	if file.Size > 5*1024*1024 {
		h.logger.Error("file too large", "size", file.Size)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Image file is too large (max 5MB)"})
		return
	}

	// Проверяем тип файла (должно быть изображение)
	contentType := file.Header.Get("Content-Type")
	if !isAllowedImageType(contentType) {
		h.logger.Error("invalid image type", "content_type", contentType)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid image type. Allowed types: image/jpeg, image/png"})
		return
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		h.logger.Error("failed to open image file", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to open image file: " + err.Error()})
		return
	}
	defer src.Close()

	// Читаем данные изображения
	imageData, err := io.ReadAll(src)
	if err != nil {
		h.logger.Error("failed to read image data", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to read image data: " + err.Error()})
		return
	}

	// Создаем страницу
	page := domain.Page{
		ChapterID: chapterID,
		Number:    number,
	}

	id, err := h.chapterService.AddPage(*c, page, imageData)
	if err != nil {
		h.logger.Error("failed to add page", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to add page: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Page added successfully",
	})
}

// deleteChapterPage удаляет страницу по ID
// @Summary Удалить страницу
// @Description Удаляет страницу по её ID
// @Tags chapters
// @Accept json
// @Produce json
// @Param page_id path int true "ID страницы"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/chapters/pages/{page_id} [delete]
func (h *ChapterHandler) deleteChapterPage(c *gin.Context) {
	pageID, err := strconv.Atoi(c.Param("page_id"))
	if err != nil {
		h.logger.Error("invalid page id format", "id", c.Param("page_id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid page ID format"})
		return
	}

	err = h.chapterService.DeletePage(*c, pageID)
	if err != nil {
		h.logger.Error("failed to delete page", "id", pageID, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to delete page: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Page deleted successfully",
	})
}

// authMiddleware middleware для проверки роли пользователя
func (h *ChapterHandler) authMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Предполагается, что JWT middleware уже добавил информацию о пользователе в контекст
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Unauthorized"})
			c.Abort()
			return
		}

		// Проверяем роль пользователя
		role := userRole.(string)
		if !hasRequiredRole(role, requiredRole) {
			c.JSON(http.StatusForbidden, ErrorResponse{Message: "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRequiredRole проверяет, имеет ли пользователь требуемую роль
func hasRequiredRole(userRole, requiredRole string) bool {
	// Администратор имеет все права
	if userRole == "admin" {
		return true
	}

	// Модератор имеет права модератора и пользователя
	if userRole == "moderator" && (requiredRole == "moderator" || requiredRole == "user") {
		return true
	}

	// Обычный пользователь имеет только права пользователя
	if userRole == "user" && requiredRole == "user" {
		return true
	}

	return false
}

// isAllowedImageType проверяет допустимый тип изображения
func isAllowedImageType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	return allowedTypes[contentType]
}
