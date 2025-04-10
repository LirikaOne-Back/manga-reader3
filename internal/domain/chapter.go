package domain

import (
	"time"
)

// Chapter представляет главу манги
type Chapter struct {
	ID        int       `json:"id"`
	MangaID   int       `json:"manga_id"`
	Number    float64   `json:"number"` // Номер главы (может быть дробным, например 1.5)
	Title     string    `json:"title"`
	PageCount int       `json:"page_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Page представляет страницу главы
type Page struct {
	ID        int    `json:"id"`
	ChapterID int    `json:"chapter_id"`
	Number    int    `json:"number"`
	ImageURL  string `json:"image_url"`
}
