package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/LirikaOne-Back/manga-reader3/internal/repository"
)

// UserService предоставляет методы для работы с пользователями
type UserService struct {
	repo   repository.UserRepository
	logger *slog.Logger
}

// NewUserService создает новый экземпляр UserService
func NewUserService(
	repo repository.UserRepository,
	logger *slog.Logger,
) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// GetByID возвращает пользователя по ID
func (s *UserService) GetByID(ctx context.Context, id int) (domain.User, error) {
	s.logger.Debug("getting user by id", "id", id)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user", "id", id, "error", err)
		return domain.User{}, err
	}

	return user, nil
}

// Update обновляет информацию о пользователе
func (s *UserService) Update(ctx context.Context, user domain.User) error {
	s.logger.Info("updating user", "id", user.ID)

	if user.ID == 0 {
		return errors.New("user id is required")
	}

	// Проверяем, существует ли пользователь
	_, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		s.logger.Error("user not found", "id", user.ID, "error", err)
		return err
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Error("failed to update user", "id", user.ID, "error", err)
		return err
	}

	s.logger.Info("user updated successfully", "id", user.ID)
	return nil
}

// Delete удаляет пользователя по ID
func (s *UserService) Delete(ctx context.Context, id int) error {
	s.logger.Info("deleting user", "id", id)

	// Проверяем, существует ли пользователь
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("user not found", "id", id, "error", err)
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete user", "id", id, "error", err)
		return err
	}

	s.logger.Info("user deleted successfully", "id", id)
	return nil
}

// AddBookmark добавляет закладку для пользователя
func (s *UserService) AddBookmark(ctx context.Context, userID, mangaID int) error {
	s.logger.Info("adding bookmark", "user_id", userID, "manga_id", mangaID)

	err := s.repo.AddBookmark(ctx, userID, mangaID)
	if err != nil {
		s.logger.Error("failed to add bookmark", "user_id", userID, "manga_id", mangaID, "error", err)
		return fmt.Errorf("failed to add bookmark: %w", err)
	}

	s.logger.Info("bookmark added successfully", "user_id", userID, "manga_id", mangaID)
	return nil
}

// RemoveBookmark удаляет закладку пользователя
func (s *UserService) RemoveBookmark(ctx context.Context, userID, mangaID int) error {
	s.logger.Info("removing bookmark", "user_id", userID, "manga_id", mangaID)

	err := s.repo.RemoveBookmark(ctx, userID, mangaID)
	if err != nil {
		s.logger.Error("failed to remove bookmark", "user_id", userID, "manga_id", mangaID, "error", err)
		return fmt.Errorf("failed to remove bookmark: %w", err)
	}

	s.logger.Info("bookmark removed successfully", "user_id", userID, "manga_id", mangaID)
	return nil
}

// GetBookmarks возвращает список закладок пользователя
func (s *UserService) GetBookmarks(ctx context.Context, userID int) ([]domain.Manga, error) {
	s.logger.Debug("getting bookmarks", "user_id", userID)

	mangas, err := s.repo.GetBookmarks(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get bookmarks", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get bookmarks: %w", err)
	}

	return mangas, nil
}

// SaveReadHistory сохраняет историю чтения
func (s *UserService) SaveReadHistory(ctx context.Context, history domain.ReadHistory) error {
	s.logger.Info("saving read history",
		"user_id", history.UserID,
		"manga_id", history.MangaID,
		"chapter_id", history.ChapterID)

	err := s.repo.SaveReadHistory(ctx, history)
	if err != nil {
		s.logger.Error("failed to save read history", "error", err)
		return fmt.Errorf("failed to save read history: %w", err)
	}

	s.logger.Info("read history saved successfully")
	return nil
}

// GetReadHistory возвращает историю чтения пользователя
func (s *UserService) GetReadHistory(ctx context.Context, userID int) ([]domain.ReadHistory, error) {
	s.logger.Debug("getting read history", "user_id", userID)

	history, err := s.repo.GetReadHistory(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get read history", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get read history: %w", err)
	}

	return history, nil
}
