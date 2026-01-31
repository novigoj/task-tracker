package initialize // пакет для подключения к БД

import (
	"fmt" // обёртка ошибок

	"gorm.io/driver/postgres" // драйвер Postgres для GORM
	"gorm.io/gorm"            // GORM
)

func New(databaseURL string) (*gorm.DB, error) { // создать подключение
	gormDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{}) // открыть Postgres
	if err != nil { // если не удалось
		return nil, fmt.Errorf("open db: %w", err) // вернуть ошибку
	}
	return gormDB, nil // вернуть *gorm.DB
}
