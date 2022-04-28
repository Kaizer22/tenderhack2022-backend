package impl

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"main/db"
	"main/logging"
	"main/model/entity"
	"main/utils"
)

const (
	PostgresOrmDebugHookEnv = "PRINT_ALL_QUERIES"

	PostgresHostEnv              = "POSTGRES_HOST"
	PostgresPortEnv              = "POSTGRES_PORT"
	PostgresDatabaseEnv          = "POSTGRES_DB"
	PostgresUsernameEnv          = "POSTGRES_USERNAME"
	PostgresPasswordEnv          = "POSTGRES_PASSWORD"
	PostgresSSLModeEnv           = "POSTGRES_SSL_MODE"
	PostgresConnectionTimeoutEnv = "POSTGRES_CONNECTION_TIMEOUT"
)

type pgOrmConnectionProvider struct {
	options *pg.Options
	pgDb    *pg.DB
}

func NewPgOrmConnectionProvider() (db.ConnectionProvider, error) {
	options := pg.Options{
		Addr:     utils.GetEnv(PostgresHostEnv, "localhost") + ":" + utils.GetEnv(PostgresPortEnv, "5432"),
		User:     utils.GetEnv(PostgresUsernameEnv, "postgres"),
		Password: utils.GetEnv(PostgresPasswordEnv, "postgres"),
		Database: utils.GetEnv(PostgresDatabaseEnv, "store"),
		//ReadTimeout:  time.Second * time.Duration(utils.GetEnvInt(PostgresConnectionTimeoutEnv, 10)),
		//WriteTimeout: time.Second * time.Duration(utils.GetEnvInt(PostgresConnectionTimeoutEnv, 10)),
	}

	dbConnection := pg.Connect(&options)
	//dbConnection.AddQueryHook(pgdebug.DebugHook{
	//	// Print all queries.
	//	Verbose: utils.GetEnvBool(PostgresOrmDebugHookEnv, false),
	//})
	logging.DebugFormat("opened a new pg orm connection - %v+", options)
	provider := pgOrmConnectionProvider{
		options: &options,
		pgDb:    dbConnection,
	}
	if err := provider.pgDb.Ping(context.Background()); err != nil {
		logging.FatalFormat("cannot establish a connection to %v+", options)
		return nil, err
	}

	return provider, nil
}

const ConnectionTestQuery = "SELECT 1"

func (p pgOrmConnectionProvider) IsConnected() (bool, error) {
	if _, err := p.pgDb.Exec(ConnectionTestQuery); err != nil {
		return false, err
	}
	return true, nil
}

func (p pgOrmConnectionProvider) Description() string {
	return "Postgres ORM"
}

func (p pgOrmConnectionProvider) Migrate(migrationPath string) error {
	err := p.createSchema()
	if err != nil {
		return err
	}
	return nil
}

func (p pgOrmConnectionProvider) Connection() interface{} {
	return p.pgDb
}

func (p pgOrmConnectionProvider) createSchema() error {
	models := []interface{}{
		(*entity.Category)(nil),
		(*entity.Product)(nil),
		(*entity.Profile)(nil),
		(*entity.Account)(nil),
		(*entity.Bet)(nil),
		(*entity.QuotationSession)(nil),
		(*entity.ProductJournal)(nil),
	}
	for _, model := range models {
		err := p.pgDb.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
