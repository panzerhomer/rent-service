package main

import (
	"avito/internal/config"
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath, migrationType string

	flag.StringVar(&migrationsPath, "path", "", "path to migrations")
	flag.StringVar(&migrationType, "direction", "up", "define direction of migrations")
	flag.Parse()

	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	storagePath := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.PostgresDB)

	if storagePath == "" {
		panic("-path is required")
	}

	if migrationsPath == "" {
		panic("-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		storagePath,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied")
}
