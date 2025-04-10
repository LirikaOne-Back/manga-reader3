package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/LirikaOne-Back/manga-reader3/internal/service"
	"github.com/gin-gonic/gin"
)

// Middleware содержит middleware функции для обработки HTTP-запросов
type Middleware struct {
	authService AuthService
	logger      *slog.Logger
}

// NewMiddleware создает новый экземпляр Middleware
func NewMiddleware(authService AuthService, logger *slog.Logger) *Middleware {
	return &Middleware{
		authService: authService,
		logger:      logger,
	}
}

// JWTAuth middleware для проверки JWT токена
func (m *Middleware) JWTAuth() gin.HandlerFunc {
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
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			m.logger.Warn("invalid token", "error", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		// Добавляем информацию о пользователе в контекст
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RoleAuth middleware для проверки роли пользователя
func (m *Middleware) RoleAuth(requiredRole string) gin.HandlerFunc {
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
		if !hasRole(role, requiredRole) {
			c.JSON(http.StatusForbidden, ErrorResponse{Message: "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRole проверяет, имеет ли пользователь требуемую роль
func hasRole(userRole, requiredRole string) bool {
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

// CORS middleware для настройки CORS
func (m *Middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Recover middleware для обработки паник
func (m *Middleware) Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Error("panic recovered", "error", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal server error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// Timeout middleware для ограничения времени выполнения запроса
func (m *Middleware) Timeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement request timeout
		c.Next()
	}
}

// RequestID middleware для добавления уникального ID запроса
func (m *Middleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement request ID
		c.Next()
	}
}

// Logger middleware для логирования запросов
func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := service.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := service.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		// Логируем в зависимости от статус-кода
		logFunc := m.logger.Info
		if statusCode >= 500 {
			logFunc = m.logger.Error
		} else if statusCode >= 400 {
			logFunc = m.logger.Warn
		}

		logFunc("HTTP request",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency,
			"error", errorMessage,
		)
	}
}

// ContentTypeJSON middleware для проверки Content-Type
func (m *Middleware) ContentTypeJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "multipart/form-data") {
				c.JSON(http.StatusUnsupportedMediaType, ErrorResponse{Message: "Content-Type must be application/json or multipart/form-data"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
