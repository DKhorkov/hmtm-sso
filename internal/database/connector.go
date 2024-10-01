package database

import (
	"database/sql"
	"fmt"
	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"
	"log/slog"

	"github.com/DKhorkov/hmtm-sso/internal/config"
	"github.com/DKhorkov/hmtm-sso/pkg/logging"

	_ "github.com/lib/pq" // Postgres driver
)

type CommonDBConnector struct {
	Connection *sql.DB
	Driver     string
	DSN        string
	Logger     *slog.Logger
}

func (connector *CommonDBConnector) Connect() error {
	if connector.Connection == nil {
		connection, err := sql.Open(connector.Driver, connector.DSN)

		if err != nil {
			return err
		}

		connector.Connection = connection
	}

	return nil
}

func (connector *CommonDBConnector) GetConnection() *sql.DB {
	if connector.Connection == nil {
		if err := connector.Connect(); err != nil {
			return nil
		}
	}

	return connector.Connection
}

func (connector *CommonDBConnector) GetTransaction() (*sql.Tx, error) {
	if connector.Connection == nil {
		return nil, &customerrors.NilDBConnectionError{}
	}

	return connector.Connection.Begin()
}

func (connector *CommonDBConnector) CloseConnection() {
	if connector.Connection == nil {
		return
	}

	if err := connector.Connection.Close(); err != nil {
		connector.Logger.Error(
			"Failed to close database connection",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
	}
}

func New(dbConfig config.DatabaseConfig, logger *slog.Logger) (*CommonDBConnector, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DatabaseName,
		dbConfig.SSLMode,
	)

	dbConnector := &CommonDBConnector{
		Driver: dbConfig.Driver,
		DSN:    dsn,
		Logger: logger,
	}

	if err := dbConnector.Connect(); err != nil {
		return nil, err
	}

	return dbConnector, nil
}
