package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/jmoiron/sqlx"
)

// ChapterRepo реализует интерфейс repository.ChapterRepository
type ChapterRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// NewChapterRepo создает новый репозиторий для работы с главами
func NewChapterRepo(db *sqlx.DB, logger *slog.Logger) *ChapterRepo {
	return &ChapterRepo{
		db:     db,
		logger: logger,
	}
}

// GetByMangaID возвращает список глав для указанной манги
func (r *ChapterRepo) GetByMangaID(ctx context.Context, mangaID int) ([]domain.Chapter, error) {
	r.logger.Debug("executing GetByMangaID chapters query", "manga_id", mangaID)

	query := `
		SELECT id, manga_id, number, title, page_count, created_at, updated_at
		FROM chapters
		WHERE manga_id = $1
		ORDER BY number
	`

	var chapters []domain.Chapter
	if err := r.db.SelectContext(ctx, &chapters, query, mangaID); err != nil {
		r.logger.Error("error selecting chapters by manga_id", "manga_id", mangaID, "error", err)
		return nil, fmt.Errorf("error selecting chapters: %w", err)
	}

	return chapters, nil
}

// GetByID возвращает главу по ID
func (r *ChapterRepo) GetByID(ctx context.Context, id int) (domain.Chapter, error) {
	r.logger.Debug("executing GetByID chapter query", "id", id)

	query := `
		SELECT id, manga_id, number, title, page_count, created_at, updated_at
		FROM chapters
		WHERE id = $1
	`

	var chapter domain.Chapter
	if err := r.db.GetContext(ctx, &chapter, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Chapter{}, fmt.Errorf("chapter with id %d not found", id)
		}
		r.logger.Error("error selecting chapter by id", "id", id, "error", err)
		return domain.Chapter{}, fmt.Errorf("error selecting chapter: %w", err)
	}

	return chapter, nil
}

// Create создает новую главу
func (r *ChapterRepo) Create(ctx context.Context, chapter domain.Chapter) (int, error) {
	r.logger.Debug("executing Create chapter query",
		"manga_id", chapter.MangaID,
		"number", chapter.Number,
		"title", chapter.Title)

	query := `
		INSERT INTO chapters (manga_id, number, title, page_count)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx, query,
		chapter.MangaID, chapter.Number, chapter.Title, chapter.PageCount,
	).Scan(&id)

	if err != nil {
		r.logger.Error("error inserting chapter", "error", err)
		return 0, fmt.Errorf("error inserting chapter: %w", err)
	}

	return id, nil
}

// Update обновляет информацию о главе
func (r *ChapterRepo) Update(ctx context.Context, chapter domain.Chapter) error {
	r.logger.Debug("executing Update chapter query",
		"id", chapter.ID,
		"title", chapter.Title,
		"number", chapter.Number)

	query := `
		UPDATE chapters SET 
			number = $1, 
			title = $2, 
			page_count = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(
		ctx, query,
		chapter.Number, chapter.Title, chapter.PageCount, chapter.ID,
	)

	if err != nil {
		r.logger.Error("error updating chapter", "id", chapter.ID, "error", err)
		return fmt.Errorf("error updating chapter: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error getting rows affected", "error", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chapter with id %d not found", chapter.ID)
	}

	return nil
}

// Delete удаляет главу по ID
func (r *ChapterRepo) Delete(ctx context.Context, id int) error {
	r.logger.Debug("executing Delete chapter query", "id", id)

	query := "DELETE FROM chapters WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("error deleting chapter", "id", id, "error", err)
		return fmt.Errorf("error deleting chapter: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error getting rows affected", "error", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chapter with id %d not found", id)
	}

	return nil
}

// GetPages возвращает список страниц для указанной главы
func (r *ChapterRepo) GetPages(ctx context.Context, chapterID int) ([]domain.Page, error) {
	r.logger.Debug("executing GetPages query", "chapter_id", chapterID)

	query := `
		SELECT id, chapter_id, number, image_url
		FROM pages
		WHERE chapter_id = $1
		ORDER BY number
	`

	var pages []domain.Page
	if err := r.db.SelectContext(ctx, &pages, query, chapterID); err != nil {
		r.logger.Error("error selecting pages", "chapter_id", chapterID, "error", err)
		return nil, fmt.Errorf("error selecting pages: %w", err)
	}

	return pages, nil
}

// AddPage добавляет новую страницу в главу
func (r *ChapterRepo) AddPage(ctx context.Context, page domain.Page) (int, error) {
	r.logger.Debug("executing AddPage query",
		"chapter_id", page.ChapterID,
		"number", page.Number)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("error starting transaction", "error", err)
		return 0, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Добавляем страницу
	query := `
		INSERT INTO pages (chapter_id, number, image_url)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int
	err = tx.QueryRowContext(
		ctx, query,
		page.ChapterID, page.Number, page.ImageURL,
	).Scan(&id)

	if err != nil {
		r.logger.Error("error inserting page", "error", err)
		return 0, fmt.Errorf("error inserting page: %w", err)
	}

	// Обновляем количество страниц в главе
	_, err = tx.ExecContext(ctx, `
		UPDATE chapters SET 
			page_count = (SELECT COUNT(*) FROM pages WHERE chapter_id = $1)
		WHERE id = $1
	`, page.ChapterID)

	if err != nil {
		r.logger.Error("error updating chapter page count", "chapter_id", page.ChapterID, "error", err)
		return 0, fmt.Errorf("error updating chapter page count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("error committing transaction", "error", err)
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return id, nil
}

// DeletePage удаляет страницу по ID
func (r *ChapterRepo) DeletePage(ctx context.Context, id int) error {
	r.logger.Debug("executing DeletePage query", "id", id)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("error starting transaction", "error", err)
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Получаем ID главы
	var chapterID int
	err = tx.QueryRowContext(ctx, "SELECT chapter_id FROM pages WHERE id = $1", id).Scan(&chapterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("page with id %d not found", id)
		}
		r.logger.Error("error getting chapter_id", "page_id", id, "error", err)
		return fmt.Errorf("error getting chapter_id: %w", err)
	}

	// Удаляем страницу
	query := "DELETE FROM pages WHERE id = $1"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("error deleting page", "id", id, "error", err)
		return fmt.Errorf("error deleting page: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error getting rows affected", "error", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("page with id %d not found", id)
	}

	// Обновляем нумерацию страниц
	_, err = tx.ExecContext(ctx, `
		WITH ranked AS (
			SELECT id, ROW_NUMBER() OVER (ORDER BY number) as new_number
			FROM pages
			WHERE chapter_id = $1
		)
		UPDATE pages p
		SET number = r.new_number
		FROM ranked r
		WHERE p.id = r.id
	`, chapterID)

	if err != nil {
		r.logger.Error("error updating page numbers", "chapter_id", chapterID, "error", err)
		return fmt.Errorf("error updating page numbers: %w", err)
	}

	// Обновляем количество страниц в главе
	_, err = tx.ExecContext(ctx, `
		UPDATE chapters SET 
			page_count = (SELECT COUNT(*) FROM pages WHERE chapter_id = $1)
		WHERE id = $1
	`, chapterID)

	if err != nil {
		r.logger.Error("error updating chapter page count", "chapter_id", chapterID, "error", err)
		return fmt.Errorf("error updating chapter page count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("error committing transaction", "error", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
