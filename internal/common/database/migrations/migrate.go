package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"smart-hub/internal/common/logger"
)

func RunMigrations(databaseURL string) error {
	return RunMigrationsWithSourceUrl(databaseURL, "file://migrations")
}

func RunMigrationsWithSourceUrl(databaseURL string, sourceUrl string) error {
	logger.Info("Running migrations...")
	m, err := migrate.New(
		sourceUrl,
		databaseURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Info("Migrations completed successfully")
	return nil
}
