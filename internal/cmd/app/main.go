package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

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

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Name,
		cfg.DatabaseConfig.SslMode,
	)

	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	if err = db.Ping(); err != nil {
		log.Printf("%v", err)
		return
	}

	migrations, err := migrate.New(
		"file://migrations"+string(os.PathSeparator),
		dbUrl,
	)

	if err != nil {
		log.Printf("Failed to create migration instance: %v", err)
		return
	}

	if err = migrations.Up(); !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Failed to up migrations: %v", err)
		return
	}
	s := server.NewServer(db)
	s.SetupRoutes(cfg)

	log.Println(http.ListenAndServe(":8080", s.Router))
}
