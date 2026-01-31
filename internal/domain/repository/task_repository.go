package repository // интерфейс репозитория

import (
	"context"                            // ctx
	"task-tracker/internal/domain/types" // модели
)

type TaskRepository interface { // контракт хранилища
	Ping(ctx context.Context) error // проверка БД

	Create(ctx context.Context, task *types.Task) error // создать
	List(ctx context.Context, done *bool, limit, offset int) ([]types.Task, error)
	GetByID(ctx context.Context, id uint) (*types.Task, error) // получить

	Update(ctx context.Context, id uint, title *string, done *bool) (*types.Task, error) // обновить частично
	Delete(ctx context.Context, id uint) error                                           // удалить
}
