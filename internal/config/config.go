package config // пакет конфигурации

import (
	"fmt"     // ошибки/формат
	"os"      // env vars
	"strings" // TrimSpace
)

type Config struct { // конфиг приложения
	Port        string // порт сервера
	DatabaseURL string // DSN БД
}

func Load() (Config, error) { // читаем env -> Config
	port, ok := os.LookupEnv("PORT")          // PORT из env
	if !ok || strings.TrimSpace(port) == "" { // обязателен
		return Config{}, fmt.Errorf("PORT is required") // ошибка
	}

	dbURL, ok := os.LookupEnv("DATABASE_URL")  // DATABASE_URL из env
	if !ok || strings.TrimSpace(dbURL) == "" { // обязателен
		return Config{}, fmt.Errorf("DATABASE_URL is required") // ошибка
	}

	return Config{ // собираем конфиг
		Port:        strings.TrimSpace(port),  // чистим пробелы
		DatabaseURL: strings.TrimSpace(dbURL), // чистим пробелы
	}, nil
}
