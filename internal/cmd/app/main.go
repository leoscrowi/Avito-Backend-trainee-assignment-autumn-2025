package main

import (
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
)

func main() {

	// TODO: доделать dbUrl, сделать конфиг
	dbUrl := "temp var"

	migrations, err := migrate.New(
		"file://migrations",
		dbUrl,
	)

	if err != nil {
		log.Println("Failed to create migration instance", err)
		os.Exit(1)
	}

	if err = migrations.Up(); !errors.Is(err, migrate.ErrNoChange) {
		log.Println("Failed to apply migrations", err)
		os.Exit(1)
	}
}
