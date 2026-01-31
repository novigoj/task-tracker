package handlers // HTTP-хендлеры

import (
	"net/http" // HTTP статусы

	"github.com/gin-gonic/gin"            // Gin контекст
	"task-tracker/internal/domain/service" // сервисный слой
)

type VersionHandler struct { // хендлер версии
	taskService *service.TaskService // зависимость сервиса
}

func NewVersionHandler(taskService *service.TaskService) *VersionHandler { // конструктор
	return &VersionHandler{taskService: taskService} // сохранить сервис
}

func (h *VersionHandler) GetVersion(c *gin.Context) { // GET /version
	c.JSON(http.StatusOK, gin.H{ // 200 + JSON
		"version": h.taskService.Version(), // версия из сервиса
	})
}
