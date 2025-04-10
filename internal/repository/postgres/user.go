package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/jmoiron/sqlx"
)

// UserRepo реализует интерфейс repository.UserRepository
type UserRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// NewUserRepo создает новый репозиторий для работы с пользователями
func NewUserRepo(db *sqlx.DB, logger *slog.Logger) *UserRepo {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

// Create создает нового пользователя
func (r *UserRepo) Create(ctx context.Context, user domain.User) (int, error) {
	r.logger.Debug("executing Create user query", "username", user.Username, "email", user.Email)

	query := `
		INSERT INTO users (username, email, password_hash, avatar_url, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx, query,
		user.Username, user.Email, user.PasswordHash, user.AvatarURL, user.Role,
	).Scan(&id)

	if err != nil {
		r.logger.Error("error inserting user", "error", err)
		return 0, fmt.Errorf("error inserting user: %w", err)
	}

	return id, nil
}

// GetByID возвращает пользователя по ID
func (r *UserRepo) GetByID(ctx context.Context, id int) (domain.User, error) {
	r.logger.Debug("executing GetByID user query", "id", id)

	query := `
		SELECT id, username, email, password_hash, avatar_url, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id %d not found", id)
		}
		r.logger.Error("error selecting user by id", "id", id, "error", err)
		return domain.User{}, fmt.Errorf("error selecting user: %w", err)
	}

	return user, nil
}

// GetByUsername возвращает пользователя по имени пользователя
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	r.logger.Debug("executing GetByUsername query", "username", username)

	query := `
		SELECT id, username, email, password_hash, avatar_url, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with username %s not found", username)
		}
		r.logger.Error("error selecting user by username", "username", username, "error", err)
		return domain.User{}, fmt.Errorf("error selecting user: %w", err)
	}

	return user, nil
}

// GetByEmail возвращает пользователя по email
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	r.logger.Debug("executing GetByEmail query", "email", email)

	query := `
		SELECT id, username, email, password_hash, avatar_url, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with email %s not found", email)
		}
		r.logger.Error("error selecting user by email", "email", email, "error", err)
		return domain.User{}, fmt.Errorf("error selecting user: %w", err)
	}

	return user, nil
}

// Update обновляет информацию о пользователе
func (r *UserRepo) Update(ctx context.Context, user domain.User) error {
	r.logger.Debug("executing Update user query", "id", user.ID, "username", user.Username)

	query := `
		UPDATE users SET 
			username = $1, 
			email = $2, 
			password_hash = $3, 
			avatar_url = $4, 
			role = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(
		ctx, query,
		user.Username, user.Email, user.PasswordHash, user.AvatarURL, user.Role, user.ID,
	)

	if err != nil {
		r.logger.Error("error updating user", "id", user.ID, "error", err)
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error getting rows affected", "error", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}

	return nil
}

// Delete удаляет пользователя по ID
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	r.logger.Debug("executing Delete user query", "id", id)

	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("error deleting user", "id", id, "error", err)
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error getting rows affected", "error", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

// AddBookmark добавляет закладку для пользователя
func (r *UserRepo) AddBookmark(ctx context.Context, userID, mangaID int) error {
	r.logger.Debug("executing AddBookmark query", "user_id", userID, "manga_id", mangaID)

	query := `
		INSERT INTO bookmarks (user_id, manga_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, manga_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, userID, mangaID)
	if err != nil {
		r.logger.Error("error adding bookmark", "user_id", userID, "manga_id", mangaID, "error", err)
		return fmt.Errorf("error adding bookmark: %w", err)
	}

	return nil
}

// RemoveBookmark удаляет закладку пользователя
func (r *UserRepo) RemoveBookmark(ctx context.Context, userID, mangaID int) error {
	r.logger.Debug("executing RemoveBookmark query", "user_id", userID, "manga_id", mangaID)

	query := "DELETE FROM bookmarks WHERE user_id = $1 AND manga_id = $2"
	result, err := r.db.ExecContext(ctx, query, userID, mangaID)
	if err != nil {
		r.logger.Error("error removing bookmark", "user_id", userID, "manga_id", mangaID, "error", err)
		return fmt.Errorf("error removing bookmark: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error getting rows affected", "error", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("bookmark not found for user %d and manga %d", userID, mangaID)
	}

	return nil
}

// GetBookmarks возвращает список закладок пользователя
func (r *UserRepo) GetBookmarks(ctx context.Context, userID int) ([]domain.Manga, error) {
	r.logger.Debug("executing GetBookmarks query", "user_id", userID)

	query := `
		SELECT m.id, m.title, m.alter_title, m.description, m.cover_url, 
		m.year, m.status, m.author, m.artist, m.rating, 
		m.created_at, m.updated_at
		FROM manga m
		JOIN bookmarks b ON m.id = b.manga_id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`

	var mangas []domain.Manga
	if err := r.db.SelectContext(ctx, &mangas, query, userID); err != nil {
		r.logger.Error("error selecting bookmarks", "user_id", userID, "error", err)
		return nil, fmt.Errorf("error selecting bookmarks: %w", err)
	}

	// Получаем жанры для каждой манги
	for i := range mangas {
		genres, err := r.getMangaGenres(ctx, mangas[i].ID)
		if err != nil {
			r.logger.Error("error getting manga genres", "manga_id", mangas[i].ID, "error", err)
			continue
		}
		mangas[i].Genres = genres
	}

	return mangas, nil
}

// SaveReadHistory сохраняет историю чтения
func (r *UserRepo) SaveReadHistory(ctx context.Context, history domain.ReadHistory) error {
	r.logger.Debug("executing SaveReadHistory query",
		"user_id", history.UserID,
		"manga_id", history.MangaID,
		"chapter_id", history.ChapterID)

	query := `
		INSERT INTO read_history (user_id, manga_id, chapter_id, page, read_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, manga_id, chapter_id) DO UPDATE
		SET page = $4, read_at = $5
	`

	readAt := history.ReadAt
	if readAt.IsZero() {
		readAt = time.Now()
	}

	_, err := r.db.ExecContext(
		ctx, query,
		history.UserID, history.MangaID, history.ChapterID, history.Page, readAt,
	)

	if err != nil {
		r.logger.Error("error saving read history", "error", err)
		return fmt.Errorf("error saving read history: %w", err)
	}

	return nil
}

// GetReadHistory возвращает историю чтения пользователя
func (r *UserRepo) GetReadHistory(ctx context.Context, userID int) ([]domain.ReadHistory, error) {
	r.logger.Debug("executing GetReadHistory query", "user_id", userID)

	query := `
		SELECT id, user_id, manga_id, chapter_id, page, read_at
		FROM read_history
		WHERE user_id = $1
		ORDER BY read_at DESC
	`

	var history []domain.ReadHistory
	if err := r.db.SelectContext(ctx, &history, query, userID); err != nil {
		r.logger.Error("error selecting read history", "user_id", userID, "error", err)
		return nil, fmt.Errorf("error selecting read history: %w", err)
	}

	return history, nil
}

// getMangaGenres возвращает жанры для указанной манги
func (r *UserRepo) getMangaGenres(ctx context.Context, mangaID int) ([]domain.Genre, error) {
	query := `
		SELECT g.id, g.name
		FROM genres g
		JOIN manga_genres mg ON g.id = mg.genre_id
		WHERE mg.manga_id = $1
		ORDER BY g.name
	`

	var genres []domain.Genre
	if err := r.db.SelectContext(ctx, &genres, query, mangaID); err != nil {
		return nil, fmt.Errorf("error selecting manga genres: %w", err)
	}

	return genres, nil
}
