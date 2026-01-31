package response // общий слой ответов

import (
	"errors"   // errors.As
	"net/http" // HTTP статусы

	"task-tracker/internal/domain/service" // ошибки сервиса

	"github.com/gin-gonic/gin" // Gin
)

type ErrorResponse struct { // JSON формат ошибки
	Error   string `json:"error"`             // код ошибки
	Details any    `json:"details,omitempty"` // детали (опц.)
}

func JSONError(c *gin.Context, status int, code string, details any) { // отдать ошибку JSON
	// регистрируем ошибку в контексте (для middleware логирования)
	// (если details — не error, то просто логируем код)
	c.JSON(status, ErrorResponse{ // ответ клиенту
		Error:   code,    // код
		Details: details, // детали
	})
}

func FromServiceError(c *gin.Context, err error) { // маппинг ошибок сервиса -> HTTP
	var appErr *service.AppError // ожидаем AppError
	if errors.As(err, &appErr) { // приводим ошибку
		switch appErr.Code { // по коду
		case service.CodeValidation: // validation
			c.Error(err)                                                             // сохранить в контексте
			JSONError(c, http.StatusBadRequest, string(appErr.Code), appErr.Details) // 400
			return
		case service.CodeNotFound: // not_found
			c.Error(err)
			JSONError(c, http.StatusNotFound, string(appErr.Code), appErr.Details) // 404
			return
		default: // всё остальное
			c.Error(err)
			JSONError(c, http.StatusInternalServerError, string(service.CodeInternal), nil) // 500
			return
		}
	}

	// Любая “неизвестная” ошибка → internal_error
	c.Error(err)                                                                    // сохранить в контексте
	JSONError(c, http.StatusInternalServerError, string(service.CodeInternal), nil) // 500
}
