package repositories_test

import (
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/DKhorkov/libs/db"
	"github.com/pressly/goose/v3"
)

const (
	driver = "sqlite3"
	// dsn    = "file::memory:?cache=shared".
	dsn              = "../../test.db"
	migrationsDir    = "/migrations"
	gooseZeroVersion = 0
)

func StartUp(t *testing.T) *db.CommonDBConnector {
	dbConnector, err := db.New(dsn, driver, &slog.Logger{})
	if err != nil {
		t.Fatal(err)
	}

	if err = goose.SetDialect(driver); err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	err = goose.Up(
		dbConnector.Pool(),
		path.Dir(
			path.Dir(cwd),
		)+migrationsDir,
	)

	if err != nil {
		panic(err)
	}

	return dbConnector
}

func TearDown(t *testing.T, dbConnector db.Connector) {
	defer func() {
		if err := dbConnector.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	err = goose.DownTo(
		dbConnector.Pool(),
		path.Dir(
			path.Dir(cwd),
		)+migrationsDir,
		gooseZeroVersion,
	)

	if err != nil {
		panic(err)
	}
}
