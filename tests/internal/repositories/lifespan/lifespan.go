package testlifespan

import (
	"database/sql"
	"os"
	"path"
	"testing"

	testconfig "github.com/DKhorkov/hmtm-sso/tests/config"
	"github.com/pressly/goose/v3"
)

var testsConfig = testconfig.New()

func StartUp(t *testing.T) *sql.DB {
	connection, err := sql.Open(testsConfig.Database.Driver, testsConfig.Database.DSN)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err = goose.SetDialect(testsConfig.Database.Driver); err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	err = goose.Up(
		connection,
		path.Dir(
			path.Dir(
				path.Dir(cwd),
			),
		)+testsConfig.Database.MigrationsDir,
	)

	if err != nil {
		panic(err)
	}

	return connection
}

func TearDown(t *testing.T, connection *sql.DB) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	err = goose.Down(
		connection,
		path.Dir(
			path.Dir(
				path.Dir(cwd),
			),
		)+testsConfig.Database.MigrationsDir,
	)

	if err != nil {
		panic(err)
	}

	if err = connection.Close(); err != nil {
		t.Fatalf("failed to close the connection to the database: %v", err)
	}
}
