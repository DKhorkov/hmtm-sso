//go:build integration

package repositories_test

import (
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/pressly/goose/v3"

	"github.com/DKhorkov/libs/db"
)

const (
	driver = "sqlite3"
	// dsn    = "file::memory:?cache=shared".
	dsn              = "../../test.db"
	migrationsDir    = "/migrations"
	gooseZeroVersion = 0
)

func StartUp(tb testing.TB) db.Connector {
	dbConnector, err := db.New(dsn, driver, &slog.Logger{})
	if err != nil {
		tb.Fatal(err)
	}

	//// Cleaning up not to use defer Teardown in test and to avoid DRY.
	defer tb.Cleanup(
		func() {
			TearDown(tb, dbConnector)
		},
	)

	if err = goose.SetDialect(driver); err != nil {
		tb.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		tb.Fatalf("failed to get cwd: %v", err)
	}

	err = goose.Up(
		dbConnector.Pool(),
		path.Dir(
			path.Dir(cwd),
		)+migrationsDir,
	)

	if err != nil {
		tb.Fatal(err)
	}

	return dbConnector
}

func TearDown(tb testing.TB, dbConnector db.Connector) {
	defer func() {
		if err := dbConnector.Close(); err != nil {
			tb.Fatal(err)
		}
	}()

	cwd, err := os.Getwd()
	if err != nil {
		tb.Fatalf("failed to get cwd: %v", err)
	}

	err = goose.DownTo(
		dbConnector.Pool(),
		path.Dir(
			path.Dir(cwd),
		)+migrationsDir,
		gooseZeroVersion,
	)

	if err != nil {
		tb.Fatal(err)
	}
}
