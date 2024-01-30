package main

import (
	"encoding/json"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
	"time"
)

type MigrationsConfig struct {
	StoragePath    string `json:"storage_path"`
	MigrationsPath string `json:"migrations_path"`
}

func main() {
	logFile, _ := os.Create("logs" + time.Now().String())
	defer logFile.Close()

	configFile, err := os.Open("cmd/migrator/config/migrations_config.json")
	if err != nil {
		logFile.WriteString(err.Error())
		panic(err.Error())
	}

	var cfg MigrationsConfig
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&cfg); err != nil {
		logFile.WriteString(err.Error())
		panic(err.Error())
	}

	logFile.WriteString("Конфигурация завершена")

	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		cfg.StoragePath,
	)
	defer m.Close()
	if err != nil {
		logFile.WriteString(err.Error())
		panic(err.Error())
	}
	logFile.WriteString("КНачало миграции")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logFile.WriteString(err.Error())
		panic(err.Error())
	}
	logFile.WriteString("Миграция завершена")
}
