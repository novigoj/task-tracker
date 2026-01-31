package repository // реализации репозиториев

import (
	"context" // ctx
	"errors"  // errors.Is

	"task-tracker/internal/domain/types" // модели

	"gorm.io/gorm" // GORM
)

type TaskGormRepository struct { // repo на GORM
	db *gorm.DB // подключение
}

func NewTaskGormRepository(db *gorm.DB) *TaskGormRepository { // конструктор
	return &TaskGormRepository{db: db} // сохранить db
}

func (r *TaskGormRepository) Ping(ctx context.Context) error { // ping БД
	sqlDB, err := r.db.DB() // получить *sql.DB
	if err != nil {
		return err // вернуть ошибку
	}
	return sqlDB.PingContext(ctx) // ping с ctx
}

func (r *TaskGormRepository) Create(ctx context.Context, task *types.Task) error { // создать задачу
	return r.db.WithContext(ctx).Create(task).Error // INSERT task
}

func (r *TaskGormRepository) List(ctx context.Context, done *bool, limit, offset int) ([]types.Task, error) { // список задач
	var tasks []types.Task // результат

	q := r.db.WithContext(ctx).Model(&types.Task{}).Order("id") // базовый запрос
	if done != nil {                                            // фильтр done?
		q = q.Where("done = ?", *done) // WHERE done=...
	}
	if limit > 0 { // лимит
		q = q.Limit(limit) // LIMIT
	}
	if offset > 0 { // сдвиг
		q = q.Offset(offset) // OFFSET
	}

	err := q.Find(&tasks).Error // выполнить SELECT
	return tasks, err           // вернуть
}

func (r *TaskGormRepository) GetByID(ctx context.Context, id uint) (*types.Task, error) { // получить по id
	var task types.Task                                 // объект
	err := r.db.WithContext(ctx).First(&task, id).Error // SELECT ... WHERE id=?
	if err != nil {                                     // обработка ошибок
		if errors.Is(err, gorm.ErrRecordNotFound) { // нет записи
			return nil, ErrNotFound // доменная not found
		}
		return nil, err // прочие ошибки
	}
	return &task, nil // вернуть задачу
}

func (r *TaskGormRepository) Update(ctx context.Context, id uint, title *string, done *bool) (*types.Task, error) { // частичный апдейт
	task, err := r.GetByID(ctx, id) // загрузить
	if err != nil {
		return nil, err // ErrNotFound уже тут
	}

	if title != nil { // менять title?
		task.Title = *title
	}
	if done != nil { // менять done?
		task.Done = *done
	}

	if err := r.db.WithContext(ctx).Save(task).Error; err != nil { // сохранить
		return nil, err
	}
	return task, nil // вернуть
}

func (r *TaskGormRepository) Delete(ctx context.Context, id uint) error { // удалить по id
	res := r.db.WithContext(ctx).Delete(&types.Task{}, id) // DELETE ... WHERE id=?
	if res.Error != nil {                                  // ошибка
		return res.Error
	}
	if res.RowsAffected == 0 { // не удалилось
		return ErrNotFound
	}
	return nil // ok
}
