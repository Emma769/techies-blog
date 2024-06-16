package pgstore

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (store *PGStore) MigrateDB() error {
	instance, err := postgres.WithInstance(store.db, &postgres.Config{})
	if err != nil {
		return err
	}

	file := "file:///migrations"

	m, err := migrate.NewWithDatabaseInstance(file, Driver_Name, instance)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
