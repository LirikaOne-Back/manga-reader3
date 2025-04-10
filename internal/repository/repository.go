package repository

import (
	"context"
	"github.com/LirikaOne-Back/manga-reader3/internal/domain"
)

// MangaRepository определяет методы для работы с мангой
type MangaRepository interface {
	GetAll(ctx context.Context, filter domain.MangaFilter) ([]domain.Manga, int, error)
	GetByID(ctx context.Context, id int) (domain.Manga, error)
	Create(ctx context.Context, manga domain.Manga) (int, error)
	Update(ctx context.Context, manga domain.Manga) error
	Delete(ctx context.Context, id int) error
	GetGenres(ctx context.Context) ([]domain.Genre, error)
}

// ChapterRepository определяет методы для работы с главами
type ChapterRepository interface {
	GetByMangaID(ctx context.Context, mangaID int) ([]domain.Chapter, error)
	GetByID(ctx context.Context, id int) (domain.Chapter, error)
	Create(ctx context.Context, chapter domain.Chapter) (int, error)
	Update(ctx context.Context, chapter domain.Chapter) error
	Delete(ctx context.Context, id int) error
	GetPages(ctx context.Context, chapterID int) ([]domain.Page, error)
	AddPage(ctx context.Context, page domain.Page) (int, error)
	DeletePage(ctx context.Context, id int) error
}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetByID(ctx context.Context, id int) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id int) error

	// Методы для работы с закладками
	AddBookmark(ctx context.Context, userID, mangaID int) error
	RemoveBookmark(ctx context.Context, userID, mangaID int) error
	GetBookmarks(ctx context.Context, userID int) ([]domain.Manga, error)

	// Методы для работы с историей чтения
	SaveReadHistory(ctx context.Context, history domain.ReadHistory) error
	GetReadHistory(ctx context.Context, userID int) ([]domain.ReadHistory, error)
}

// Repositories объединяет все репозитории для удобного внедрения зависимостей
type Repositories struct {
	Manga   MangaRepository
	Chapter ChapterRepository
	User    UserRepository
}
