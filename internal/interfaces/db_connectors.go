package interfaces

import "database/sql"

type DBConnector interface {
	Connect() error
	CloseConnection()
	GetTransaction() (*sql.Tx, error)
	GetConnection() *sql.DB
}
