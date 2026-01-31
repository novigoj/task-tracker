package types // пакет с моделями/типами

import "time" // time.Time

type User struct { // модель пользователя (GORM)
	ID        uint      `gorm:"primaryKey"`            // PK
	Email     string    `gorm:"uniqueIndex;not null"`  // уникальный email, обязателен
	CreatedAt time.Time                                // автозаполняется GORM

	Tasks []Task `gorm:"foreignKey:UserID"` // связь 1->many по UserID
}
