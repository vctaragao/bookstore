package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func Migrate(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/postgres",
	}

	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}
