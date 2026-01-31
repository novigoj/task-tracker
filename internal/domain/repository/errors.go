package repository // пакет репозиториев

import "errors" // errors.New

var ErrNotFound = errors.New("not found") // общая ошибка "не найдено"
