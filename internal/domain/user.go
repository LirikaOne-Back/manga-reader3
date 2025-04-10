package domain

import (
	"time"
)

// User представляет пользователя системы
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Не выводим в JSON
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Role         string    `json:"role"` // user, moderator, admin
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Bookmark представляет закладку пользователя
type Bookmark struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	MangaID   int       `json:"manga_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ReadHistory представляет историю чтения
type ReadHistory struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	MangaID   int       `json:"manga_id"`
	ChapterID int       `json:"chapter_id"`
	Page      int       `json:"page"`
	ReadAt    time.Time `json:"read_at"`
}

// UserSignup представляет данные для регистрации пользователя
type UserSignup struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserLogin представляет данные для входа
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse представляет ответ с токеном доступа
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
