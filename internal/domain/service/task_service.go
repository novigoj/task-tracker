package service // сервисный слой

import (
	"context" // ctx
	"errors"  // errors.Is
	"strings" // TrimSpace

	"task-tracker/internal/domain/repository" // repo интерфейс + ошибки
	"task-tracker/internal/domain/types"      // модели
)

type TaskService struct { // сервис задач
	repo repository.TaskRepository // зависимость
}

func NewTaskService(repo repository.TaskRepository) *TaskService { // конструктор
	return &TaskService{repo: repo} // сохранить repo
}

func (s *TaskService) Version() string { return "0.1.0" } // версия

func (s *TaskService) Health(ctx context.Context) error { // healthcheck
	return s.repo.Ping(ctx) // ping хранилища
}

func (s *TaskService) Create(ctx context.Context, userID uint, title string) (*types.Task, error) { // создать задачу
	title = strings.TrimSpace(title) // чистим title
	if userID == 0 || title == "" {  // базовая валидация
		return nil, Validation(map[string]string{
			"user_id": "must be > 0",
			"title":   "required",
		}) // ошибка валидации
	}

	task := &types.Task{ // собираем модель
		UserID: userID, // владелец
		Title:  title,  // заголовок
		Done:   false,  // дефолт
	}

	if err := s.repo.Create(ctx, task); err != nil { // записываем в БД
		return nil, Internal(err) // пробрасываем ошибку
	}
	return task, nil // вернуть созданную
}

func (s *TaskService) List(ctx context.Context, done *bool, limit, offset int) ([]types.Task, error) { // список задач
	// базовые правила для API
	if limit <= 0 { // дефолт
		limit = 20
	}
	if limit > 100 || offset < 0 {
		return nil, Validation(map[string]string{
			"limit":  "must be 1..100",
			"offset": "must be >= 0",
		})
	}
	tasks, err := s.repo.List(ctx, done, limit, offset)
	if err != nil {
		return nil, Internal(err)
	}
	return tasks, nil
}

func (s *TaskService) GetByID(ctx context.Context, id uint) (*types.Task, error) { // получить задачу
	task, err := s.repo.GetByID(ctx, id) // repo вызов
	if err != nil {                      // маппим ошибки
		if errors.Is(err, repository.ErrNotFound) { // нет записи
			return nil, NotFound(nil) // ошибка сервиса
		}
		return nil, Internal(err) // прочее
	}
	return task, nil // ok
}

func (s *TaskService) Update(ctx context.Context, id uint, title *string, done *bool) (*types.Task, error) { // PATCH задачи
	if id == 0 { // id обязателен
		return nil, Validation(map[string]string{"id": "required"})
	}
	if title == nil && done == nil { // нечего менять
		return nil, Validation(map[string]string{
			"title": "required",
			"done":  "required",
		})
	}

	if title != nil { // валидируем title
		t := strings.TrimSpace(*title) // trim
		if t == "" {                   // пусто нельзя
			return nil, Validation(map[string]string{"title": "required"})
		}
		title = &t // подменяем на очищенный
	}

	task, err := s.repo.Update(ctx, id, title, done) // обновление в repo
	if err != nil {                                  // маппим ошибки
		if errors.Is(err, repository.ErrNotFound) { // нет записи
			return nil, NotFound(nil)
		}
		return nil, Internal(err) // прочее
	}
	return task, nil // ok
}

func (s *TaskService) Delete(ctx context.Context, id uint) error { // удалить задачу
	if id == 0 { // id обязателен
		return Validation(map[string]string{"id": "required"})
	}
	err := s.repo.Delete(ctx, id) // удалить в repo
	if err != nil {               // маппим ошибки
		if errors.Is(err, repository.ErrNotFound) { // не найдено
			return NotFound(nil)
		}
		return Internal(err) // прочее
	}
	return nil // ok
}
