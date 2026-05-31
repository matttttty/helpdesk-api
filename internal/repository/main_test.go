package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("helpdesk_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		fmt.Printf("failed to start container: %v\n", err)
		os.Exit(1)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Printf("failed to get dsn: %v\n", err)
		os.Exit(1)
	}

	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("failed to open db: %v\n", err)
		os.Exit(1)
	}

	if err := runMigration(testDB); err != nil {
		fmt.Printf("failed to migrate: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	_ = testDB.Close()
	_ = pgContainer.Terminate(ctx)
	os.Exit(code)
}

func runMigration(db *sql.DB) error {
	sqlBytes, err := os.ReadFile("../../migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	_, err = db.Exec(string(sqlBytes))
	return err
}
