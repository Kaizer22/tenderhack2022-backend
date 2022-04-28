package db

const (
	PgMigrationsPath = "file://db/migration/postgres"
)

type ConnectionProvider interface {
	Connection() interface{}
	Description() string
	IsConnected() (bool, error)
	Migrate(migrationPath string) error
}
