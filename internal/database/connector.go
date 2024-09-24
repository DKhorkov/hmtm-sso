package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/hmtm-sso/internal/config"
	"github.com/DKhorkov/hmtm-sso/pkg/logging"

	_ "github.com/lib/pq" // Postgres driver
)

type CommonDBConnector struct {
	connection *sql.DB
	driver     string
	dsn        string
	logger     *slog.Logger
}

func (connector *CommonDBConnector) Connect() error {
	if connector.connection == nil {
		connection, err := sql.Open(connector.driver, connector.dsn)

		if err != nil {
			return err
		}

		connector.connection = connection
	}

	return nil
}

func (connector *CommonDBConnector) GetConnection() *sql.DB {
	return connector.connection
}

func (connector *CommonDBConnector) GetTransaction() (*sql.Tx, error) {
	return connector.connection.Begin()
}

func (connector *CommonDBConnector) CloseConnection() {
	if connector.connection == nil {
		return
	}

	if err := connector.connection.Close(); err != nil {
		connector.logger.Error(
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
		driver: dbConfig.Driver,
		dsn:    dsn,
		logger: logger,
	}

	if err := dbConnector.Connect(); err != nil {
		return nil, err
	}

	return dbConnector, nil
}
