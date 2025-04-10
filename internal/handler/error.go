package handler

// ErrorResponse стандартный формат ответа с ошибкой
type ErrorResponse struct {
	Message string `json:"message"`
}
