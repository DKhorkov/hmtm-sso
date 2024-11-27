package repositories_test

import (
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/DKhorkov/libs/db"
	"github.com/pressly/goose/v3"
)

var testsConfig = db.NewTestConfig()

func StartUp(t *testing.T) *db.CommonDBConnector {
	dbConnector, err := db.NewTestConnector(testsConfig, &slog.Logger{})
	if err != nil {
		t.Fatal(err)
	}

	if err = goose.SetDialect(testsConfig.Driver); err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	err = goose.Up(
		dbConnector.GetConnection(),
		path.Dir(
			path.Dir(cwd),
		)+testsConfig.MigrationsDir,
	)

	if err != nil {
		panic(err)
	}

	return dbConnector
}

func TearDown(t *testing.T, dbConnector db.Connector) {
	defer dbConnector.CloseConnection()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	err = goose.Down(
		dbConnector.GetConnection(),
		path.Dir(
			path.Dir(cwd),
		)+testsConfig.MigrationsDir,
	)

	if err != nil {
		panic(err)
	}
}
