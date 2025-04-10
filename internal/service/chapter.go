package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/LirikaOne-Back/manga-reader3/internal/repository"
)

// ChapterService предоставляет методы для работы с главами
type ChapterService struct {
	repo       repository.ChapterRepository
	mangaRepo  repository.MangaRepository
	logger     *slog.Logger
	imagesPath string
}

// NewChapterService создает новый экземпляр ChapterService
func NewChapterService(
	repo repository.ChapterRepository,
	mangaRepo repository.MangaRepository,
	logger *slog.Logger,
	imagesPath string,
) *ChapterService {
	return &ChapterService{
		repo:       repo,
		mangaRepo:  mangaRepo,
		logger:     logger,
		imagesPath: imagesPath,
	}
}

// GetByMangaID возвращает список глав для указанной манги
func (s *ChapterService) GetByMangaID(ctx context.Context, mangaID int) ([]domain.Chapter, error) {
	s.logger.Debug("getting chapters by manga id", "manga_id", mangaID)

	// Проверяем, существует ли манга
	_, err := s.mangaRepo.GetByID(ctx, mangaID)
	if err != nil {
		s.logger.Error("manga not found", "manga_id", mangaID, "error", err)
		return nil, fmt.Errorf("manga with id %d not found: %w", mangaID, err)
	}

	chapters, err := s.repo.GetByMangaID(ctx, mangaID)
	if err != nil {
		s.logger.Error("failed to get chapters", "manga_id", mangaID, "error", err)
		return nil, err
	}

	return chapters, nil
}

// GetByID возвращает главу по ID
func (s *ChapterService) GetByID(ctx context.Context, id int) (domain.Chapter, error) {
	s.logger.Debug("getting chapter by id", "id", id)

	chapter, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get chapter", "id", id, "error", err)
		return domain.Chapter{}, err
	}

	return chapter, nil
}

// Create создает новую главу
func (s *ChapterService) Create(ctx context.Context, chapter domain.Chapter) (int, error) {
	s.logger.Info("creating new chapter", "manga_id", chapter.MangaID, "number", chapter.Number)

	if chapter.MangaID == 0 {
		return 0, errors.New("manga id is required")
	}

	if chapter.Number <= 0 {
		return 0, errors.New("chapter number must be positive")
	}

	if chapter.Title == "" {
		return 0, errors.New("chapter title is required")
	}

	// Проверяем, существует ли манга
	_, err := s.mangaRepo.GetByID(ctx, chapter.MangaID)
	if err != nil {
		s.logger.Error("manga not found", "manga_id", chapter.MangaID, "error", err)
		return 0, fmt.Errorf("manga with id %d not found: %w", chapter.MangaID, err)
	}

	// Создаем каталог для изображений главы
	err = s.createChapterImageDir(chapter.MangaID, chapter.Number)
	if err != nil {
		s.logger.Error("failed to create chapter image directory", "error", err)
		return 0, fmt.Errorf("failed to create chapter image directory: %w", err)
	}

	id, err := s.repo.Create(ctx, chapter)
	if err != nil {
		s.logger.Error("failed to create chapter", "error", err)
		return 0, err
	}

	s.logger.Info("chapter created successfully", "id", id)
	return id, nil
}

// Update обновляет информацию о главе
func (s *ChapterService) Update(ctx context.Context, chapter domain.Chapter) error {
	s.logger.Info("updating chapter", "id", chapter.ID)

	if chapter.ID == 0 {
		return errors.New("chapter id is required")
	}

	// Проверяем, существует ли глава
	existingChapter, err := s.repo.GetByID(ctx, chapter.ID)
	if err != nil {
		s.logger.Error("chapter not found", "id", chapter.ID, "error", err)
		return err
	}

	// Если номер главы изменился, обновляем каталог для изображений
	if existingChapter.Number != chapter.Number {
		err = s.updateChapterImageDir(existingChapter.MangaID, existingChapter.Number, chapter.Number)
		if err != nil {
			s.logger.Error("failed to update chapter image directory", "error", err)
			return fmt.Errorf("failed to update chapter image directory: %w", err)
		}
	}

	err = s.repo.Update(ctx, chapter)
	if err != nil {
		s.logger.Error("failed to update chapter", "id", chapter.ID, "error", err)
		return err
	}

	s.logger.Info("chapter updated successfully", "id", chapter.ID)
	return nil
}

// Delete удаляет главу по ID
func (s *ChapterService) Delete(ctx context.Context, id int) error {
	s.logger.Info("deleting chapter", "id", id)

	// Получаем информацию о главе перед удалением
	chapter, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("chapter not found", "id", id, "error", err)
		return err
	}

	// Удаляем главу из БД
	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete chapter", "id", id, "error", err)
		return err
	}

	// Удаляем каталог с изображениями главы
	err = s.deleteChapterImageDir(chapter.MangaID, chapter.Number)
	if err != nil {
		s.logger.Warn("failed to delete chapter image directory", "id", id, "error", err)
		// Не возвращаем ошибку, так как глава уже удалена из БД
	}

	s.logger.Info("chapter deleted successfully", "id", id)
	return nil
}

// GetPages возвращает список страниц для указанной главы
func (s *ChapterService) GetPages(ctx context.Context, chapterID int) ([]domain.Page, error) {
	s.logger.Debug("getting pages for chapter", "chapter_id", chapterID)

	// Проверяем, существует ли глава
	_, err := s.repo.GetByID(ctx, chapterID)
	if err != nil {
		s.logger.Error("chapter not found", "chapter_id", chapterID, "error", err)
		return nil, err
	}

	pages, err := s.repo.GetPages(ctx, chapterID)
	if err != nil {
		s.logger.Error("failed to get pages", "chapter_id", chapterID, "error", err)
		return nil, err
	}

	return pages, nil
}

// AddPage добавляет новую страницу в главу
func (s *ChapterService) AddPage(ctx context.Context, page domain.Page, imageData []byte) (int, error) {
	s.logger.Info("adding page to chapter", "chapter_id", page.ChapterID, "number", page.Number)

	if page.ChapterID == 0 {
		return 0, errors.New("chapter id is required")
	}

	// Проверяем, существует ли глава
	chapter, err := s.repo.GetByID(ctx, page.ChapterID)
	if err != nil {
		s.logger.Error("chapter not found", "chapter_id", page.ChapterID, "error", err)
		return 0, err
	}

	// Проверяем и устанавливаем номер страницы, если не указан
	if page.Number <= 0 {
		// Получаем существующие страницы
		pages, err := s.repo.GetPages(ctx, page.ChapterID)
		if err != nil {
			s.logger.Error("failed to get pages", "chapter_id", page.ChapterID, "error", err)
			return 0, err
		}
		page.Number = len(pages) + 1
	}

	// Сохраняем изображение
	imagePath, err := s.savePageImage(chapter.MangaID, chapter.Number, page.Number, imageData)
	if err != nil {
		s.logger.Error("failed to save page image", "error", err)
		return 0, fmt.Errorf("failed to save page image: %w", err)
	}

	// Устанавливаем URL изображения
	page.ImageURL = imagePath

	id, err := s.repo.AddPage(ctx, page)
	if err != nil {
		s.logger.Error("failed to add page", "error", err)
		// Удаляем изображение, если не удалось добавить страницу
		_ = os.Remove(filepath.Join(s.imagesPath, imagePath))
		return 0, err
	}

	s.logger.Info("page added successfully", "id", id, "chapter_id", page.ChapterID)
	return id, nil
}

// DeletePage удаляет страницу по ID
func (s *ChapterService) DeletePage(ctx context.Context, id int) error {
	s.logger.Info("deleting page", "id", id)

	// Получаем информацию о странице перед удалением
	pages, err := s.repo.GetPages(ctx, 0) // 0 - заглушка, так как GetPages не использует chapterID
	if err != nil {
		s.logger.Error("failed to get pages", "error", err)
		return err
	}

	var targetPage domain.Page
	for _, p := range pages {
		if p.ID == id {
			targetPage = p
			break
		}
	}

	if targetPage.ID == 0 {
		return fmt.Errorf("page with id %d not found", id)
	}

	// Получаем информацию о главе
	chapter, err := s.repo.GetByID(ctx, targetPage.ChapterID)
	if err != nil {
		s.logger.Error("failed to get chapter", "chapter_id", targetPage.ChapterID, "error", err)
		return err
	}

	// Удаляем страницу из БД
	err = s.repo.DeletePage(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete page", "id", id, "error", err)
		return err
	}

	// Удаляем изображение
	err = s.deletePageImage(chapter.MangaID, chapter.Number, targetPage.Number)
	if err != nil {
		s.logger.Warn("failed to delete page image", "id", id, "error", err)
		// Не возвращаем ошибку, так как страница уже удалена из БД
	}

	s.logger.Info("page deleted successfully", "id", id)
	return nil
}

// createChapterImageDir создает каталог для изображений главы
func (s *ChapterService) createChapterImageDir(mangaID int, chapterNumber float64) error {
	dirPath := filepath.Join(s.imagesPath, fmt.Sprintf("manga_%d/chapter_%.2f", mangaID, chapterNumber))
	return os.MkdirAll(dirPath, 0755)
}

// updateChapterImageDir обновляет каталог для изображений главы при изменении номера главы
func (s *ChapterService) updateChapterImageDir(mangaID int, oldNumber, newNumber float64) error {
	oldPath := filepath.Join(s.imagesPath, fmt.Sprintf("manga_%d/chapter_%.2f", mangaID, oldNumber))
	newPath := filepath.Join(s.imagesPath, fmt.Sprintf("manga_%d/chapter_%.2f", mangaID, newNumber))

	// Проверяем, существует ли старый каталог
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		// Если старого каталога нет, просто создаем новый
		return os.MkdirAll(newPath, 0755)
	}

	// Переименовываем каталог
	return os.Rename(oldPath, newPath)
}

// deleteChapterImageDir удаляет каталог с изображениями главы
func (s *ChapterService) deleteChapterImageDir(mangaID int, chapterNumber float64) error {
	dirPath := filepath.Join(s.imagesPath, fmt.Sprintf("manga_%d/chapter_%.2f", mangaID, chapterNumber))
	return os.RemoveAll(dirPath)
}

// savePageImage сохраняет изображение страницы и возвращает относительный путь
func (s *ChapterService) savePageImage(mangaID int, chapterNumber float64, pageNumber int, imageData []byte) (string, error) {
	// Формируем относительный путь к изображению
	relPath := fmt.Sprintf("manga_%d/chapter_%.2f/page_%03d.jpg", mangaID, chapterNumber, pageNumber)
	absPath := filepath.Join(s.imagesPath, relPath)

	// Создаем каталог, если он не существует
	err := os.MkdirAll(filepath.Dir(absPath), 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Сохраняем изображение
	err = os.WriteFile(absPath, imageData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write image file: %w", err)
	}

	return relPath, nil
}

// deletePageImage удаляет изображение страницы
func (s *ChapterService) deletePageImage(mangaID int, chapterNumber float64, pageNumber int) error {
	relPath := fmt.Sprintf("manga_%d/chapter_%.2f/page_%03d.jpg", mangaID, chapterNumber, pageNumber)
	absPath := filepath.Join(s.imagesPath, relPath)
	return os.Remove(absPath)
}
