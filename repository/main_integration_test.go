package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	dsn := os.Getenv("DB_TEST_DSN")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to test DB: %v", err)
	}

	testDB = pool
	defer pool.Close()

	code := m.Run()
	os.Exit(code)
}
