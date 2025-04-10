package main

import (
	"log"
	"os"

	"github.com/LirikaOne-Back/manga-reader3/internal/app"
	"github.com/LirikaOne-Back/manga-reader3/internal/config"
	"github.com/LirikaOne-Back/manga-reader3/pkg/logger"

	_ "github.com/LirikaOne-Back/manga-reader3/docs" // Import для swagger
)

// @title Manga Reader API
// @version 1.0
// @description API сервер для читалки манги

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализируем логгер
	log := logger.New(cfg.Logger)
	log.Info("Starting manga reader API server")

	// Создаем каталог для изображений, если он не существует
	if err := os.MkdirAll(cfg.Storage.ImagesPath, 0755); err != nil {
		log.Error("Failed to create images directory", "error", err)
		os.Exit(1)
	}

	// Инициализируем приложение
	app, err := app.NewApp(cfg, log)
	if err != nil {
		log.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Запускаем приложение
	if err := app.Run(); err != nil {
		log.Error("Application error", "error", err)
		os.Exit(1)
	}
}
