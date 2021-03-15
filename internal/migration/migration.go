package migration

import (
	"errors"
	"net/url"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigratePostgres(sourceURL, databaseURL string) error {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return err
	}

	databaseURL = u.String()

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return err
	}

	return nil
}
