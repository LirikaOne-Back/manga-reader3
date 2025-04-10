package domain

import (
	"time"
)

// Genre представляет жанр манги
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Manga представляет информацию о манге
type Manga struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	AlterTitle  string    `json:"alter_title,omitempty"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url"`
	Year        int       `json:"year"`
	Status      string    `json:"status"` // ongoing, completed, hiatus
	Author      string    `json:"author"`
	Artist      string    `json:"artist,omitempty"`
	Rating      float64   `json:"rating"`
	Genres      []Genre   `json:"genres"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MangaFilter содержит параметры для фильтрации списка манги
type MangaFilter struct {
	Genre    string
	Status   string
	SortBy   string // title, rating, date
	SortDesc bool
	Search   string
	Page     int
	PageSize int
}
