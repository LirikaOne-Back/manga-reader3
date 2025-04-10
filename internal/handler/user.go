package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/gin-gonic/gin"
)

// UserHandler обрабатывает HTTP-запросы, связанные с пользователями
type UserHandler struct {
	userService UserService
	logger      *slog.Logger
	middleware  *Middleware
}

// UserService интерфейс сервиса пользователей
type UserService interface {
	GetByID(ctx gin.Context, id int) (domain.User, error)
	Update(ctx gin.Context, user domain.User) error
	Delete(ctx gin.Context, id int) error
	AddBookmark(ctx gin.Context, userID, mangaID int) error
	RemoveBookmark(ctx gin.Context, userID, mangaID int) error
	GetBookmarks(ctx gin.Context, userID int) ([]domain.Manga, error)
	SaveReadHistory(ctx gin.Context, history domain.ReadHistory) error
	GetReadHistory(ctx gin.Context, userID int) ([]domain.ReadHistory, error)
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userService UserService, middleware *Middleware, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		middleware:  middleware,
		logger:      logger,
	}
}

// Register регистрирует обработчики путей для пользователей
func (h *UserHandler) Register(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		// Пути, требующие аутентификации
		authenticated := users.Group("/")
		authenticated.Use(h.middleware.JWTAuth())
		{
			authenticated.GET("/profile", h.getUserProfile)
			authenticated.PUT("/profile", h.updateUserProfile)
			authenticated.DELETE("/profile", h.deleteUserProfile)

			// Пути для работы с закладками
			bookmarks := authenticated.Group("/bookmarks")
			{
				bookmarks.GET("", h.getUserBookmarks)
				bookmarks.POST("/:manga_id", h.addUserBookmark)
				bookmarks.DELETE("/:manga_id", h.removeUserBookmark)
			}

			// Пути для работы с историей чтения
			history := authenticated.Group("/history")
			{
				history.GET("", h.getUserReadHistory)
				history.POST("", h.saveUserReadHistory)
			}
		}

		// Пути, требующие прав администратора
		admin := users.Group("/")
		admin.Use(h.middleware.JWTAuth(), h.middleware.RoleAuth("admin"))
		{
			admin.GET("/:id", h.getUserByID)
			admin.PUT("/:id", h.updateUser)
			admin.DELETE("/:id", h.deleteUser)
		}
	}
}

// getUserProfile возвращает профиль текущего пользователя
// @Summary Получить профиль пользователя
// @Description Возвращает профиль текущего аутентифицированного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/profile [get]
func (h *UserHandler) getUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int)

	user, err := h.userService.GetByID(*c, id)
	if err != nil {
		h.logger.Error("failed to get user profile", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get user profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// updateUserProfile обновляет профиль текущего пользователя
// @Summary Обновить профиль пользователя
// @Description Обновляет профиль текущего аутентифицированного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "Данные для обновления профиля"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/profile [put]
func (h *UserHandler) updateUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int)

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("invalid user data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid user data: " + err.Error()})
		return
	}

	// Устанавливаем ID из контекста
	user.ID = id

	// Проверяем, не пытается ли пользователь изменить свою роль
	currentUser, err := h.userService.GetByID(*c, id)
	if err != nil {
		h.logger.Error("failed to get current user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get current user: " + err.Error()})
		return
	}

	// Сохраняем текущую роль
	user.Role = currentUser.Role

	err = h.userService.Update(*c, user)
	if err != nil {
		h.logger.Error("failed to update user profile", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to update user profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
	})
}

// deleteUserProfile удаляет профиль текущего пользователя
// @Summary Удалить профиль пользователя
// @Description Удаляет профиль текущего аутентифицированного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/profile [delete]
func (h *UserHandler) deleteUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int)

	err := h.userService.Delete(*c, id)
	if err != nil {
		h.logger.Error("failed to delete user profile", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to delete user profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile deleted successfully",
	})
}

// getUserBookmarks возвращает закладки текущего пользователя
// @Summary Получить закладки пользователя
// @Description Возвращает список закладок текущего аутентифицированного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} domain.Manga
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/bookmarks [get]
func (h *UserHandler) getUserBookmarks(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int)

	bookmarks, err := h.userService.GetBookmarks(*c, id)
	if err != nil {
		h.logger.Error("failed to get user bookmarks", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get user bookmarks: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookmarks)
}

// addUserBookmark добавляет закладку для текущего пользователя
// @Summary Добавить закладку
// @Description Добавляет закладку для текущего аутентифицированного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param manga_id path int true "ID манги"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/bookmarks/{manga_id} [post]
func (h *UserHandler) addUserBookmark(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int)

	mangaID, err := strconv.Atoi(c.Param("manga_id"))
	if err != nil {
		h.logger.Error("invalid manga_id format", "manga_id", c.Param("manga_id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid manga ID format"})
		return
	}

	err = h.userService.AddBookmark(*c, id, mangaID)
	if err != nil {
		h.logger.Error("failed to add bookmark", "user_id", id, "manga_id", mangaID, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to add bookmark: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bookmark added successfully",
	})
}

// removeUserBookmark удаляет закладку пользователя
// @Summary Удалить закладку
// @Description Удаляет закладку текущего аутентифицированного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param manga_id path int true "ID манги"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/bookmarks/{manga_id} [delete]
func (h *UserHandler) removeUserBookmark(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int)

	mangaID, err := strconv.Atoi(c.Param("manga_id"))
	if err != nil {
		h.logger.Error("invalid manga_id format", "manga_id", c.Param("manga_id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid manga ID format"})
		return
	}

	err = h.userService.RemoveBookmark(*c, id, mangaID)
	if err != nil {
		h.logger.Error("failed to remove bookmark", "user_id", id, "manga_id", mangaID, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to remove bookmark: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bookmark removed successfully",
	})
}

// getUserReadHistory возвращает историю чтения пользователя
// @Summary Получить историю чтения
// @Description Возвращает историю чтения текущего аутентифицированного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} domain.ReadHistory
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/history [get]
func (h *UserHandler) getUserReadHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int)

	history, err := h.userService.GetReadHistory(*c, id)
	if err != nil {
		h.logger.Error("failed to get read history", "user_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get read history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// saveUserReadHistory сохраняет историю чтения пользователя
// @Summary Сохранить историю чтения
// @Description Сохраняет историю чтения текущего аутентифицированного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param history body domain.ReadHistory true "Данные истории чтения"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/history [post]
func (h *UserHandler) saveUserReadHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int)

	var history domain.ReadHistory
	if err := c.ShouldBindJSON(&history); err != nil {
		h.logger.Error("invalid history data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid history data: " + err.Error()})
		return
	}

	// Устанавливаем ID пользователя из контекста
	history.UserID = id
	history.ReadAt = Now()

	err := h.userService.SaveReadHistory(*c, history)
	if err != nil {
		h.logger.Error("failed to save read history", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to save read history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Read history saved successfully",
	})
}

// getUserByID возвращает пользователя по ID (только для администраторов)
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по ID (требуются права администратора)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/{id} [get]
func (h *UserHandler) getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid user id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid user ID format"})
		return
	}

	user, err := h.userService.GetByID(*c, id)
	if err != nil {
		h.logger.Error("failed to get user", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "User not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// updateUser обновляет пользователя по ID (только для администраторов)
// @Summary Обновить пользователя
// @Description Обновляет пользователя по ID (требуются права администратора)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param user body domain.User true "Данные для обновления пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/{id} [put]
func (h *UserHandler) updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid user id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid user ID format"})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("invalid user data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid user data: " + err.Error()})
		return
	}

	// Устанавливаем ID из URL
	user.ID = id

	err = h.userService.Update(*c, user)
	if err != nil {
		h.logger.Error("failed to update user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to update user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

// deleteUser удаляет пользователя по ID (только для администраторов)
// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID (требуются права администратора)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/users/{id} [delete]
func (h *UserHandler) deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("invalid user id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid user ID format"})
		return
	}

	err = h.userService.Delete(*c, id)
	if err != nil {
		h.logger.Error("failed to delete user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// Now возвращает текущее время (для удобства мокирования в тестах)
func Now() time.Time {
	return service.Now()
}
