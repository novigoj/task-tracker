package middleware // middleware слой

import (
	"log"

	"github.com/gin-gonic/gin"
) // логирование

// Gin

func ErrorLogger() gin.HandlerFunc { // лог ошибок из контекста
	return func(c *gin.Context) { // middleware
		c.Next() // выполнить хендлеры

		if len(c.Errors) > 0 { // есть ошибки?
			last := c.Errors.Last()                                        // берём последнюю
			log.Printf("request error method=%s path=%s status=%d err=%v", // формат лога
				c.Request.Method,   // метод
				c.Request.URL.Path, // путь
				c.Writer.Status(),  // статус ответа
				last.Err,           // ошибка
			)
		}
	}
}
