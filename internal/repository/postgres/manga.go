package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
	"github.com/jmoiron/sqlx"
)

// MangaRepo реализует интерфейс repository.MangaRepository
type MangaRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// NewMangaRepo создает новый репозиторий для работы с мангой
func NewMangaRepo(db *sqlx.DB, logger *slog.Logger) *MangaRepo {
	return &MangaRepo{
		db:     db,
		logger: logger,
	}
}

// GetAll возвращает список манги с фильтрацией и пагинацией
func (r *MangaRepo) GetAll(ctx context.Context, filter domain.MangaFilter) ([]domain.Manga, int, error) {
	r.logger.Debug("executing GetAll manga query with filter",
		"genre", filter.Genre,
		"status", filter.Status,
		"search", filter.Search)

	// Формируем запрос для получения манги
	query := `
		SELECT m.id, m.title, m.alter_title, m.description, m.cover_url, 
		m.year, m.status, m.author, m.artist, m.rating, 
		m.created_at, m.updated_at
		FROM manga m
	`
	// Формируем запрос для подсчета общего количества
	countQuery := `SELECT COUNT(*) FROM manga m`

	// Массив для хранения параметров запроса
	args := []interface{}{}
	// Массив для хранения условий WHERE
	conditions := []string{}

	// Добавляем условие для поиска по жанру
	if filter.Genre != "" {
		query = query + " JOIN manga_genres mg ON m.id = mg.manga_id JOIN genres g ON mg.genre_id = g.id"
		countQuery = countQuery + " JOIN manga_genres mg ON m.id = mg.manga_id JOIN genres g ON mg.genre_id = g.id"
		conditions = append(conditions, "g.name = ?")
		args = append(args, filter.Genre)
	}

	// Добавляем условие для фильтрации по статусу
	if filter.Status != "" {
		conditions = append(conditions, "m.status = ?")
		args = append(args, filter.Status)
	}

	// Добавляем условие для поиска по тексту
	if filter.Search != "" {
		conditions = append(conditions, "(m.title ILIKE ? OR m.alter_title ILIKE ?)")
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Добавляем все условия в запрос
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// Добавляем сортировку
	switch filter.SortBy {
	case "title":
		query += " ORDER BY m.title"
	case "rating":
		query += " ORDER BY m.rating"
	case "date":
		query += " ORDER BY m.created_at"
	default:
		query += " ORDER BY m.updated_at"
	}

	// Порядок сортировки
	if filter.SortDesc {
		query += " DESC"
	} else {
		query += " ASC"
	}

	// Пагинация
	query += " LIMIT ? OFFSET ?"
	limit := filter.PageSize
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, limit, offset)

	// Заменяем ? на $1, $2 и т.д. для PostgreSQL
	query = r.replaceQuestionMark(query)
	countQuery = r.replaceQuestionMark(countQuery)

	// Получаем общее количество
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args[:len(args)-2]...); err != nil {
		r.logger.Error("error counting total manga", "error", err)
		return nil, 0, fmt.Errorf("error counting manga: %w", err)
	}

	// Если общее количество 0, возвращаем пустой слайс
	if total == 0 {
		return []domain.Manga{}, 0, nil
	}

	// Получаем список манги
	var mangas []domain.Manga
	if err := r.db.SelectContext(ctx, &mangas, query, args...); err != nil {
		r.logger.Error("error selecting manga", "error", err)
		return nil, 0, fmt.Errorf("error selecting manga: %w", err)
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

	return mangas, total, nil
}

// GetByID возвращает мангу по ID
func (r *MangaRepo) GetByID(ctx context.Context, id int) (domain.Manga, error) {
	r.logger.Debug("executing GetByID manga query", "id", id)

	query := `
		SELECT id, title, alter_title, description, cover_url, 
		year, status, author, artist, rating, 
		created_at, updated_at
		FROM manga
		WHERE id = $1
	`

	var manga domain.Manga
	if err := r.db.GetContext(ctx, &manga, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Manga{}, fmt.Errorf("manga with id %d not found", id)
		}
		r.logger.Error("error selecting manga by id", "id", id, "error", err)
		return domain.Manga{}, fmt.Errorf("error selecting manga: %w", err)
	}

	// Получаем жанры для манги
	genres, err := r.getMangaGenres(ctx, manga.ID)
	if err != nil {
		r.logger.Error("error getting manga genres", "manga_id", manga.ID, "error", err)
	} else {
		manga.Genres = genres
	}

	return manga, nil
}

// Create создает новую мангу
func (r *MangaRepo) Create(ctx context.Context, manga domain.Manga) (int, error) {
	r.logger.Debug("executing Create manga query", "title", manga.Title)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("error starting transaction", "error", err)
		return 0, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO manga (
			title, alter_title, description, cover_url, 
			year, status, author, artist, rating
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id
	`

	var id int
	err = tx.QueryRowContext(
		ctx, query,
		manga.Title, manga.AlterTitle, manga.Description, manga.CoverURL,
		manga.Year, manga.Status, manga.Author, manga.Artist, manga.Rating,
	).Scan(&id)

	if err != nil {
		r.logger.Error("error inserting manga", "error", err)
		return 0, fmt.Errorf("error inserting manga: %w", err)
	}

	// Добавляем жанры манги
	if len(manga.Genres) > 0 {
		if err := r.insertMangaGenres(ctx, tx, id, manga.Genres); err != nil {
			r.logger.Error("error inserting manga genres", "manga_id", id, "error", err)
			return 0, fmt.Errorf("error inserting manga genres: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("error committing transaction", "error", err)
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return id, nil
}

// Update обновляет информацию о манге
func (r *MangaRepo) Update(ctx context.Context, manga domain.Manga) error {
	r.logger.Debug("executing Update manga query", "id", manga.ID, "title", manga.Title)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("error starting transaction", "error", err)
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE manga SET 
			title = $1, 
			alter_title = $2, 
			description = $3, 
			cover_url = $4, 
			year = $5, 
			status = $6, 
			author = $7, 
			artist = $8, 
			rating = $9
		WHERE id = $10
	`

	result, err := tx.ExecContext(
		ctx, query,
		manga.Title, manga.AlterTitle, manga.Description, manga.CoverURL,
		manga.Year, manga.Status, manga.Author, manga.Artist, manga.Rating,
		manga.ID,
	)

	if err != nil {
		r.logger.Error("error updating manga", "id", manga.ID, "error", err)
		return fmt.Errorf("error updating manga: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error getting rows affected", "error", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("manga with id %d not found", manga.ID)
	}

	// Обновляем жанры манги
	if len(manga.Genres) > 0 {
		// Удаляем существующие жанры
		_, err = tx.ExecContext(ctx, "DELETE FROM manga_genres WHERE manga_id = $1", manga.ID)
		if err != nil {
			r.logger.Error("error deleting manga genres", "manga_id", manga.ID, "error", err)
			return fmt.Errorf("error deleting manga genres: %w", err)
		}

		// Добавляем новые жанры
		if err := r.insertMangaGenres(ctx, tx, manga.ID, manga.Genres); err != nil {
			r.logger.Error("error inserting manga genres", "manga_id", manga.ID, "error", err)
			return fmt.Errorf("error inserting manga genres: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("error committing transaction", "error", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// Delete удаляет мангу по ID
func (r *MangaRepo) Delete(ctx context.Context, id int) error {
	r.logger.Debug("executing Delete manga query", "id", id)

	query := "DELETE FROM manga WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("error deleting manga", "id", id, "error", err)
		return fmt.Errorf("error deleting manga: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error getting rows affected", "error", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("manga with id %d not found", id)
	}

	return nil
}

// GetGenres возвращает список всех жанров
func (r *MangaRepo) GetGenres(ctx context.Context) ([]domain.Genre, error) {
	r.logger.Debug("executing GetGenres query")

	query := "SELECT id, name FROM genres ORDER BY name"
	var genres []domain.Genre
	if err := r.db.SelectContext(ctx, &genres, query); err != nil {
		r.logger.Error("error selecting genres", "error", err)
		return nil, fmt.Errorf("error selecting genres: %w", err)
	}

	return genres, nil
}

// getMangaGenres возвращает жанры для указанной манги
func (r *MangaRepo) getMangaGenres(ctx context.Context, mangaID int) ([]domain.Genre, error) {
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

// insertMangaGenres добавляет жанры для манги
func (r *MangaRepo) insertMangaGenres(ctx context.Context, tx *sql.Tx, mangaID int, genres []domain.Genre) error {
	for _, genre := range genres {
		query := "INSERT INTO manga_genres (manga_id, genre_id) VALUES ($1, $2)"
		_, err := tx.ExecContext(ctx, query, mangaID, genre.ID)
		if err != nil {
			return fmt.Errorf("error inserting manga genre: %w", err)
		}
	}
	return nil
}

// replaceQuestionMark заменяет ? на $1, $2 и т.д. для PostgreSQL
func (r *MangaRepo) replaceQuestionMark(query string) string {
	paramIndex := 0
	for strings.Contains(query, "?") {
		paramIndex++
		query = strings.Replace(query, "?", fmt.Sprintf("$%d", paramIndex), 1)
	}
	return query
}
