package handlers // HTTP-хендлеры

import (
	// errors.Is
	"net/http" // HTTP статусы
	"strconv"  // parse id

	"github.com/gin-gonic/gin" // Gin

	"task-tracker/internal/api/rest/dto" // DTO
	"task-tracker/internal/api/rest/response"
	"task-tracker/internal/domain/service" // сервис
)

type TaskHandler struct { // хендлер задач
	taskService *service.TaskService // зависимость
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler { // конструктор
	return &TaskHandler{taskService: taskService} // сохранить сервис
}

func toTaskResponse(t any) dto.TaskResponse { // маппер (пока заглушка)
	// будем приводить в местах вызова, чтобы проще читалось
	return dto.TaskResponse{}
}

func (h *TaskHandler) Create(c *gin.Context) { // POST /tasks
	var req dto.CreateTaskRequest                  // тело запроса
	if err := c.ShouldBindJSON(&req); err != nil { // распарсить JSON
		response.JSONError(c, http.StatusBadRequest, "validation_error", map[string]string{"json": "invalid"})
		return
	}

	task, err := h.taskService.Create(c.Request.Context(), req.UserID, req.Title) // создать задачу
	if err != nil {                                                               // обработка ошибок
		response.FromServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.TaskResponse{ // 201 + DTO
		ID:        task.ID,        // id
		UserID:    task.UserID,    // user
		Title:     task.Title,     // title
		Done:      task.Done,      // done
		CreatedAt: task.CreatedAt, // created
	})
}

func (h *TaskHandler) List(c *gin.Context) { // GET /tasks
	// done (optional)
	var donePtr *bool                              // nil = без фильтра
	if doneStr := c.Query("done"); doneStr != "" { // ?done=true/false
		v, err := strconv.ParseBool(doneStr) // парсим bool
		if err != nil {                      // не bool
			response.JSONError(c, http.StatusBadRequest, "validation_error", map[string]string{"done": "invalid"})
			return
		}
		donePtr = &v // включаем фильтр
	}

	// limit/offset (optional)
	limit := 20 // дефолт
	offset := 0 // дефолт

	if s := c.Query("limit"); s != "" { // ?limit=...
		v, err := strconv.Atoi(s) // парсим int
		if err != nil || v <= 0 { // не число / <=0
			response.JSONError(c, http.StatusBadRequest, "validation_error", map[string]string{"limit": "invalid"})
			return
		}
		limit = v // применяем
	}

	if s := c.Query("offset"); s != "" { // ?offset=...
		v, err := strconv.Atoi(s) // парсим int
		if err != nil || v < 0 {  // не число / <0
			response.JSONError(c, http.StatusBadRequest, "validation_error", map[string]string{"offset": "invalid"})
			return
		}
		offset = v // применяем
	}

	tasks, err := h.taskService.List(c.Request.Context(), donePtr, limit, offset) // вызов сервиса
	if err != nil {                                                               // обработка ошибок
		response.FromServiceError(c, err)
		return
	}

	resp := make([]dto.TaskResponse, 0, len(tasks)) // DTO список
	for _, t := range tasks {                       // маппинг в DTO
		resp = append(resp, dto.TaskResponse{
			ID:        t.ID,
			UserID:    t.UserID,
			Title:     t.Title,
			Done:      t.Done,
			CreatedAt: t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp) // 200 + список
}

func (h *TaskHandler) GetByID(c *gin.Context) { // GET /tasks/:id
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 { // не число / 0
		response.JSONError(c, http.StatusBadRequest, "validation_error", map[string]string{"id": "must be positive integer"})
		return
	}

	task, err := h.taskService.GetByID(c.Request.Context(), uint(id64)) // получить задачу
	if err != nil {                                                     // обработка ошибок
		response.FromServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TaskResponse{ // 200 + DTO
		ID:        task.ID,
		UserID:    task.UserID,
		Title:     task.Title,
		Done:      task.Done,
		CreatedAt: task.CreatedAt,
	})
}

func (h *TaskHandler) Update(c *gin.Context) { // PATCH /tasks/:id
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 { // не число / 0
		response.JSONError(c, http.StatusBadRequest, "validation_error", map[string]string{"id": "invalid"})
		return
	}

	var req dto.UpdateTaskRequest                  // тело PATCH
	if err := c.ShouldBindJSON(&req); err != nil { // парсим JSON
		response.JSONError(c, http.StatusBadRequest, "validation_error", map[string]string{"json": "invalid"})
		return
	}

	task, err := h.taskService.Update(c.Request.Context(), uint(id64), req.Title, req.Done) // вызов сервиса
	if err != nil {                                                                         // обработка ошибок
		response.FromServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TaskResponse{ // 200 + DTO
		ID:        task.ID,
		UserID:    task.UserID,
		Title:     task.Title,
		Done:      task.Done,
		CreatedAt: task.CreatedAt,
	})
}

func (h *TaskHandler) Delete(c *gin.Context) { // DELETE /tasks/:id
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 { // не число / 0
		response.JSONError(c, http.StatusBadRequest, "validation_error", map[string]string{"id": "invalid"})
		return
	}

	if err := h.taskService.Delete(c.Request.Context(), uint(id64)); err != nil { // удалить через сервис
		response.FromServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent) // 204 без тела
}
