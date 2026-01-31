package service // ошибки сервисного слоя

import "errors" // errors.New

var (
	ErrValidation = errors.New("validation error") // ошибка валидации
	ErrNotFound   = errors.New("not found")        // сущность не найдена
)
