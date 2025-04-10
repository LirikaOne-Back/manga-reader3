package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/LirikaOne-Back/manga-reader3/internal/repository"
)

// MangaService предоставляет методы для работы с мангой
type MangaService struct {
	repo   repository.MangaRepository
	logger *slog.Logger
}

// NewMangaService создает новый экземпляр MangaService
func NewMangaService(repo repository.MangaRepository, logger *slog.Logger) *MangaService {
	return &MangaService{
		repo:   repo,
		logger: logger,
	}
}

// GetAll возвращает список манги с фильтрацией и пагинацией
func (s *MangaService) GetAll(ctx context.Context, filter domain.MangaFilter) ([]domain.Manga, int, error) {
	s.logger.Debug("getting manga list with filter",
		"genre", filter.Genre,
		"status", filter.Status,
		"search", filter.Search,
		"page", filter.Page)

	// Устанавливаем значения по умолчанию для пагинации
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	// Получаем данные из репозитория
	mangas, total, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get manga list", "error", err)
		return nil, 0, err
	}

	s.logger.Debug("got manga list", "count", len(mangas), "total", total)
	return mangas, total, nil
}

// GetByID возвращает мангу по ID
func (s *MangaService) GetByID(ctx context.Context, id int) (domain.Manga, error) {
	s.logger.Debug("getting manga by id", "id", id)

	manga, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get manga by id", "id", id, "error", err)
		return domain.Manga{}, err
	}

	return manga, nil
}

// Create создает новую мангу
func (s *MangaService) Create(ctx context.Context, manga domain.Manga) (int, error) {
	s.logger.Info("creating new manga", "title", manga.Title)

	if manga.Title == "" {
		return 0, errors.New("manga title is required")
	}

	id, err := s.repo.Create(ctx, manga)
	if err != nil {
		s.logger.Error("failed to create manga", "title", manga.Title, "error", err)
		return 0, err
	}

	s.logger.Info("manga created successfully", "id", id, "title", manga.Title)
	return id, nil
}

// Update обновляет информацию о манге
func (s *MangaService) Update(ctx context.Context, manga domain.Manga) error {
	s.logger.Info("updating manga", "id", manga.ID, "title", manga.Title)

	if manga.ID == 0 {
		return errors.New("manga id is required")
	}

	// Проверяем, существует ли манга
	_, err := s.repo.GetByID(ctx, manga.ID)
	if err != nil {
		s.logger.Error("manga not found for update", "id", manga.ID, "error", err)
		return err
	}

	err = s.repo.Update(ctx, manga)
	if err != nil {
		s.logger.Error("failed to update manga", "id", manga.ID, "error", err)
		return err
	}

	s.logger.Info("manga updated successfully", "id", manga.ID)
	return nil
}

// Delete удаляет мангу по ID
func (s *MangaService) Delete(ctx context.Context, id int) error {
	s.logger.Info("deleting manga", "id", id)

	// Проверяем, существует ли манга
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("manga not found for deletion", "id", id, "error", err)
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete manga", "id", id, "error", err)
		return err
	}

	s.logger.Info("manga deleted successfully", "id", id)
	return nil
}

// GetGenres возвращает список всех жанров
func (s *MangaService) GetGenres(ctx context.Context) ([]domain.Genre, error) {
	s.logger.Debug("getting genres list")

	genres, err := s.repo.GetGenres(ctx)
	if err != nil {
		s.logger.Error("failed to get genres list", "error", err)
		return nil, err
	}

	return genres, nil
}
