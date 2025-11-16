package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/leoscrowi/pr-assignment-service/tests/helpers"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	if helpers.TestURL == "" {
		helpers.TestURL = "http://localhost:8080"
	}

	helpers.DbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		"postgres",
		"postgres",
		"localhost",
		"5433",
		"postgres",
		"disable",
	)

	db, err := sqlx.Open("postgres", helpers.DbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v\n", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error pinging to db: %v\n", err)
	}

	err = helpers.CleanDB(db)
	if err != nil {
		log.Fatalf("error while cleaning db")
	}

	code := m.Run()

	err = helpers.CleanDB(db)
	if err != nil {
		log.Fatalf("error while cleaning db")
	}

	os.Exit(code)
}
