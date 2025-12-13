package database

import (
	"errors"
	"fmt"
	"log"

	"backend/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg *config.DatabaseConfig, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		cfg.URL(),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func RollbackMigration(cfg *config.DatabaseConfig, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		cfg.URL(),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}

func MigrateToVersion(cfg *config.DatabaseConfig, migrationsPath string, version uint) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		cfg.URL(),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Migrate(version); err != nil {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	log.Printf("Migrated to version %d successfully", version)
	return nil
}
