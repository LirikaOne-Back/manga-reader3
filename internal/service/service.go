package service

import (
	"time"
)

// Services объединяет все сервисы приложения
type Services struct {
	Manga   *MangaService
	Chapter *ChapterService
	Auth    *AuthService
	User    *UserService
}

// Now возвращает текущее время (для удобства мокирования в тестах)
func Now() time.Time {
	return time.Now()
}
