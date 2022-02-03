package databasetest

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	// pgx.
	_ "github.com/jackc/pgx/v4/stdlib"
	// libpq.
	_ "github.com/lib/pq"
)

// GetTestingDatabaseAddress return the testing database address.
func GetTestingDatabaseAddress(t *testing.T) string {
	t.Helper()

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?timezone=UTC&sslmode=disable",
		os.Getenv("TESTING_DATABASE_USER"),
		os.Getenv("TESTING_DATABASE_PASSWORD"),
		os.Getenv("TESTING_DATABASE_HOST"),
		os.Getenv("TESTING_DATABASE_PORT"),
		os.Getenv("TESTING_DATABASE_DB"),
	)
}

// RunWithDatabase runs the test using a real database if provided
// on env vars, otherwise we just skip.
// We cannot run multiples test at the same time, use -p 1 to limit
// the execution:
// go test -v -p 1 ./...
func RunWithDatabase(t *testing.T, driver string, fn func(db *sql.DB, purgeDB func())) {
	t.Helper()

	if os.Getenv("TESTING_DATABASE_HOST") == "" {
		t.SkipNow()
	}
	address := GetTestingDatabaseAddress(t)

	db, err := sql.Open(driver, address)
	if err != nil {
		t.Fatalf("could not create database connection, %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("could not connect to the local database, %v", err)
	}

	cleanupFunc := func() {
		_, err := db.Exec(`
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
`)
		if err != nil {
			t.Fatalf("error while purging database, %v", err)
		}
	}

	defer func() {
		_ = db.Close()
	}()

	// run the test
	fn(db, cleanupFunc)
}
