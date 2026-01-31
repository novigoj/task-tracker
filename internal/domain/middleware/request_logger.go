package middleware // middleware слой

import (
	"log"  // логирование
	"time" // тайминги

	"github.com/gin-gonic/gin" // Gin
)

func RequestLogger() gin.HandlerFunc { // лог запросов
	return func(c *gin.Context) { // middleware
		start := time.Now()             // старт
		method := c.Request.Method      // метод
		path := c.Request.URL.Path      // raw путь

		c.Next() // выполнить хендлеры

		status := c.Writer.Status()     // статус
		latency := time.Since(start)    // время

		route := c.FullPath() // шаблон роута
		if route == "" {      // если нет
			route = path // fallback
		}

		switch { // уровень по статусу
		case status >= 500:
			log.Printf("ERROR request method=%s path=%s status=%d latency=%s", method, route, status, latency)
		case status >= 400:
			log.Printf("WARN  request method=%s path=%s status=%d latency=%s", method, route, status, latency)
		default:
			log.Printf("INFO  request method=%s path=%s status=%d latency=%s", method, route, status, latency)
		}
	}
}
