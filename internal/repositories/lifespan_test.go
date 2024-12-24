package repositories_test

import (
	"fmt"
	"log/slog"
	"os"
	"path"

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

func StartUp() *db.CommonDBConnector {
	dbConnector, err := db.New(dsn, driver, &slog.Logger{})
	if err != nil {
		panic(err)
	}

	if err = goose.SetDialect(driver); err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get cwd: %v", err))
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

func TearDown(dbConnector db.Connector) {
	defer func() {
		if err := dbConnector.Close(); err != nil {
			panic(err)
		}
	}()

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get cwd: %v", err))
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
