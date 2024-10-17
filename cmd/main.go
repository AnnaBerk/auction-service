package main

import (
	"auction/internal/app"
	"github.com/labstack/gommon/log"
	"log/slog"
)

func main() {
	cfg, err := app.MustLoad()
	if err != nil {
		log.Fatal("Ошибка инициализации конфига", slog.Any("error", err))
	}

	dbc, err := app.InitDB(*cfg)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных", slog.Any("error", err))
	}
	// Применение миграций
	if err := app.RunMigrations(dbc); err != nil {
		log.Fatal("Ошибка выполнения миграций", slog.Any("error", err))
	}
	a, err := app.NewApp(*cfg, dbc)
	if err != nil {
		log.Fatal(err)
	}
	a.Run()
}
