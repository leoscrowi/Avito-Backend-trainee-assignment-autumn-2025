package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/leoscrowi/pr-assignment-service/internal/config"
	"github.com/leoscrowi/pr-assignment-service/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()
	log.Println("[APPLICATION]: loaded config")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Name,
		cfg.DatabaseConfig.SslMode,
	)

	log.Printf("[APPLICATION]: Connecting to database, dbURL: %s", dbUrl)
	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("[APPLICATION]: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("[APPLICATION]: %v", err)
	}
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	log.Printf("[APPLICATION]: Starting migrations")
	// TODO: заменить потом на сепаратор для пути, т.к. для винды он \\
	migrations, err := migrate.New(
		"file://migrations/",
		dbUrl,
	)

	if err != nil {
		log.Fatalf("[APPLICATION]: Failed to create migration instance: %v", err)
	}

	if err = migrations.Up(); !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("[APPLICATION]: Failed to up migrations: %v", err)
	}
	log.Printf("[APPLICATION]: Migrations successfully created")

	s := server.NewServer(db)
	s.SetupRoutes()

	log.Fatal(http.ListenAndServe(":6060", s.Router))
}
