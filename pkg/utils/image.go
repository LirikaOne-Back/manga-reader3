package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// ImageInfo содержит информацию об изображении
type ImageInfo struct {
	Width       int
	Height      int
	ContentType string
	Format      string
	Size        int
}

// ProcessImageOptions опции для обработки изображения
type ProcessImageOptions struct {
	MaxWidth     int  // Максимальная ширина
	MaxHeight    int  // Максимальная высота
	Quality      int  // Качество (0-100)
	ConvertToJPG bool // Конвертировать в JPEG
	Watermark    bool // Добавить водяной знак
}

// DefaultProcessImageOptions возвращает опции по умолчанию
func DefaultProcessImageOptions() ProcessImageOptions {
	return ProcessImageOptions{
		MaxWidth:     1920,
		MaxHeight:    2400,
		Quality:      85,
		ConvertToJPG: true,
		Watermark:    false,
	}
}

// GetImageInfo возвращает информацию об изображении
func GetImageInfo(data []byte) (*ImageInfo, error) {
	// Определяем тип изображения
	contentType := http.DetectContentType(data)

	// Декодируем изображение для получения размеров
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()

	// Определяем формат
	format := ""
	switch {
	case strings.Contains(contentType, "jpeg"):
		format = "jpeg"
	case strings.Contains(contentType, "png"):
		format = "png"
	case strings.Contains(contentType, "webp"):
		format = "webp"
	default:
		format = "unknown"
	}

	return &ImageInfo{
		Width:       bounds.Dx(),
		Height:      bounds.Dy(),
		ContentType: contentType,
		Format:      format,
		Size:        len(data),
	}, nil
}

// ProcessImage обрабатывает изображение согласно опциям
func ProcessImage(data []byte, options ProcessImageOptions) ([]byte, error) {
	// Определяем тип изображения
	contentType := http.DetectContentType(data)

	// Декодируем изображение
	var img image.Image
	var err error

	reader := bytes.NewReader(data)

	if strings.Contains(contentType, "jpeg") {
		img, err = jpeg.Decode(reader)
	} else if strings.Contains(contentType, "png") {
		img, err = png.Decode(reader)
	} else if strings.Contains(contentType, "webp") {
		img, err = webp.Decode(reader)
	} else {
		return nil, fmt.Errorf("unsupported image format: %s", contentType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Изменяем размер, если нужно
	img = resizeImage(img, options.MaxWidth, options.MaxHeight)

	// Добавляем водяной знак, если нужно
	if options.Watermark {
		img, err = addWatermark(img)
		if err != nil {
			return nil, fmt.Errorf("failed to add watermark: %w", err)
		}
	}

	// Кодируем изображение
	var buf bytes.Buffer

	if options.ConvertToJPG {
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: options.Quality})
	} else if strings.Contains(contentType, "jpeg") {
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: options.Quality})
	} else if strings.Contains(contentType, "png") {
		err = png.Encode(&buf, img)
	} else if strings.Contains(contentType, "webp") {
		err = webp.Encode(&buf, img, &webp.Options{Quality: float32(options.Quality)})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// resizeImage изменяет размер изображения, сохраняя пропорции
func resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Если изображение уже меньше максимальных размеров, ничего не делаем
	if width <= maxWidth && height <= maxHeight {
		return img
	}

	// Вычисляем новые размеры, сохраняя пропорции
	ratio := float64(width) / float64(height)

	var newWidth, newHeight int

	if width > height {
		// Ширина больше высоты
		newWidth = maxWidth
		newHeight = int(float64(newWidth) / ratio)

		if newHeight > maxHeight {
			newHeight = maxHeight
			newWidth = int(float64(newHeight) * ratio)
		}
	} else {
		// Высота больше ширины
		newHeight = maxHeight
		newWidth = int(float64(newHeight) * ratio)

		if newWidth > maxWidth {
			newWidth = maxWidth
			newHeight = int(float64(newWidth) / ratio)
		}
	}

	// Изменяем размер с высоким качеством
	return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
}

// addWatermark добавляет водяной знак на изображение
func addWatermark(img image.Image) (image.Image, error) {
	// TODO: Реализовать добавление водяного знака
	// Для простоты просто возвращаем исходное изображение
	return img, nil
}

// IsImageValid проверяет, является ли файл допустимым изображением
func IsImageValid(reader io.Reader) (bool, error) {
	// Читаем первые 512 байт для определения типа файла
	buffer := make([]byte, 512)
	_, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Определяем тип контента
	contentType := http.DetectContentType(buffer)

	// Проверяем, является ли файл изображением
	return strings.HasPrefix(contentType, "image/"), nil
}
