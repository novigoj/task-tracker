package types // пакет с моделями/типами

import "time" // time.Time

type Task struct { // модель задачи (GORM)
	ID        uint      `gorm:"primaryKey"`             // PK
	UserID    uint      `gorm:"index;not null"`         // FK на пользователя + индекс
	Title     string    `gorm:"not null"`               // заголовок обязателен
	Done      bool      `gorm:"not null;default:false"` // флаг выполнения
	CreatedAt time.Time // автозаполняется GORM
}
