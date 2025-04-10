package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/LirikaOne-Back/manga-reader3/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler обрабатывает HTTP-запросы, связанные с аутентификацией
type AuthHandler struct {
	authService AuthService
	logger      *slog.Logger
}

// AuthService интерфейс сервиса аутентификации
type AuthService interface {
	Register(ctx gin.Context, input domain.UserSignup) (domain.User, error)
	Login(ctx gin.Context, input domain.UserLogin) (domain.TokenResponse, error)
	RefreshToken(ctx gin.Context, refreshToken string) (domain.TokenResponse, error)
	Logout(ctx gin.Context, userID int) error
	ChangePassword(ctx gin.Context, userID int, oldPassword, newPassword string) error
	ValidateToken(tokenString string) (*service.JWTClaims, error)
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authService AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register регистрирует обработчики путей для аутентификации
func (h *AuthHandler) Register(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signup)
		auth.POST("/login", h.login)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.authMiddleware(), h.logout)
		auth.PUT("/password", h.authMiddleware(), h.changePassword)
		auth.GET("/me", h.authMiddleware(), h.getMe)
	}
}

// signup обработчик для регистрации нового пользователя
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя с указанными данными
// @Tags auth
// @Accept json
// @Produce json
// @Param input body domain.UserSignup true "Регистрационные данные"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/signup [post]
func (h *AuthHandler) signup(c *gin.Context) {
	var input domain.UserSignup
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("invalid registration data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid registration data: " + err.Error()})
		return
	}

	user, err := h.authService.Register(*c, input)
	if err != nil {
		h.logger.Error("failed to register user", "error", err)

		// Проверяем тип ошибки для более информативного ответа
		if strings.Contains(err.Error(), "username already exists") {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Username already exists"})
			return
		}
		if strings.Contains(err.Error(), "email already exists") {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Email already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to register user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"message":  "User registered successfully",
	})
}

// login обработчик для входа пользователя
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает токены
// @Tags auth
// @Accept json
// @Produce json
// @Param input body domain.UserLogin true "Данные для входа"
// @Success 200 {object} domain.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) login(c *gin.Context) {
	var input domain.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("invalid login data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid login data: " + err.Error()})
		return
	}

	tokenResponse, err := h.authService.Login(*c, input)
	if err != nil {
		h.logger.Error("login failed", "error", err)
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// refresh обработчик для обновления токенов
// @Summary Обновление токенов
// @Description Обновляет токены по refresh-токену
// @Tags auth
// @Accept json
// @Produce json
// @Param input body map[string]string true "Refresh токен в формате {\"refresh_token\": \"token\"}"
// @Success 200 {object} domain.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("invalid refresh token data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid refresh token data: " + err.Error()})
		return
	}

	tokenResponse, err := h.authService.RefreshToken(*c, input.RefreshToken)
	if err != nil {
		h.logger.Error("token refresh failed", "error", err)
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// logout обработчик для выхода пользователя
// @Summary Выход пользователя
// @Description Выход пользователя (инвалидация всех токенов)
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/auth/logout [post]
func (h *AuthHandler) logout(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Unauthorized"})
		return
	}

	err := h.authService.Logout(*c, userID)
	if err != nil {
		h.logger.Error("logout failed", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to logout: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// changePassword обработчик для изменения пароля
// @Summary Изменение пароля
// @Description Изменяет пароль пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body map[string]string true "Старый и новый пароли в формате {\"old_password\": \"old\", \"new_password\": \"new\"}"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/auth/password [put]
func (h *AuthHandler) changePassword(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Unauthorized"})
		return
	}

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("invalid password change data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid password data: " + err.Error()})
		return
	}

	err := h.authService.ChangePassword(*c, userID, input.OldPassword, input.NewPassword)
	if err != nil {
		h.logger.Error("password change failed", "error", err)

		if strings.Contains(err.Error(), "invalid old password") {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid old password"})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to change password: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// getMe обработчик для получения информации о текущем пользователе
// @Summary Информация о текущем пользователе
// @Description Возвращает информацию о текущем аутентифицированном пользователе
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/auth/me [get]
func (h *AuthHandler) getMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// authMiddleware middleware для проверки аутентификации
func (h *AuthHandler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Authorization header is required"})
			c.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := headerParts[1]
		claims, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			h.logger.Warn("invalid token", "error", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		// Добавляем информацию о пользователе в контекст
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		// Также можно добавить полную информацию о пользователе, если это требуется
		// user, err := h.userService.GetByID(*c, claims.UserID)
		// if err == nil {
		//     c.Set("user", user)
		// }

		c.Next()
	}
}

// getUserIDFromContext возвращает ID пользователя из контекста
func getUserIDFromContext(c *gin.Context) int {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int)
}
