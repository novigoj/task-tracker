package middleware // middleware слой

import (
	"log"      // логирование
	"net/http" // HTTP статусы

	"github.com/gin-gonic/gin" // Gin
)

func RecoveryJSON() gin.HandlerFunc { // recovery middleware
	return gin.CustomRecovery(func(c *gin.Context, recovered any) { // ловим panic
		log.Printf("panic recovered: %v", recovered) // лог паники
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{ // 500 + JSON
			"error": "internal_error", // код ошибки
		})
	})
}
