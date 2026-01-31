package dto // DTO для API

import "time" // time.Time

type CreateTaskRequest struct { // тело запроса на создание задачи
	UserID uint   `json:"user_id"` // владелец
	Title  string `json:"title"`   // заголовок
}

type TaskResponse struct { // DTO ответа задачи
	ID        uint      `json:"id"`         // id
	UserID    uint      `json:"user_id"`    // владелец
	Title     string    `json:"title"`      // заголовок
	Done      bool      `json:"done"`       // статус
	CreatedAt time.Time `json:"created_at"` // дата создания
}

type UpdateTaskRequest struct { // PATCH payload
	Title *string `json:"title,omitempty"` // менять title (если есть)
	Done  *bool   `json:"done,omitempty"`  // менять done (если есть)
}
