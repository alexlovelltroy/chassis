package chassis

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// Our Sensible defaults are set in code here and can be overridden by instances
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		ConnectionString:      "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
		MaxOpenConnections:    5,
		MaxIdleConnections:    5,
		ConnectionMaxLifetime: 300,
		ConnectionMaxIdleTime: 300,
	}
}

// ConnectDB connects to the database and returns a pointer to the connection
// object.  It is the caller's responsibility to close the connection.
// The flavor parameter is the name of the database driver to use, e.g. "postgres".
// The uri parameter is the connection string to use, e.g. "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable".
func ConnectDB(flavor, uri string) (*sqlx.DB, error) {
	db, err := sqlx.Open(flavor, uri)
	if err != nil {
		log.Fatalf("Unable to connect to the database using %s: %v", uri, err)
	}

	// Validate the connection with Ping
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping the database using %s: %v", uri, err)
	}
	return db, nil
}

// ApplyMigrationsUp applies the database migrations for Postgres Databases up to the latest version
func ApplyMigrationsUp(db *sql.DB, migrationPath string, level uint) error {
	var driver database.Driver
	var err error
	var m *migrate.Migrate

	err = db.Ping()
	if err != nil {
		log.Fatal("Couldn't ping the database: ", err)
	}
	log.Debug("Connected to the database")
	driver, err = postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err = migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver)
	if err != nil {
		return err
	}
	log.Debug("Applying migrations up to level ", level)

	if err = m.Migrate(level); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
