package testconfig

type TestDatabaseConfig struct {
	Driver        string
	DSN           string
	MigrationsDir string
}

type TestConfig struct {
	Database TestDatabaseConfig
}

func New() *TestConfig {
	return &TestConfig{
		Database: TestDatabaseConfig{
			Driver:        "sqlite3",
			DSN:           "file::memory:?cache=shared", // "test.db" can be also used
			MigrationsDir: "/internal/database/migrations",
		},
	}
}
