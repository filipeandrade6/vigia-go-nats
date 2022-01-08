package migration

import (
	"context"
	"embed"
	"fmt"

	"github.com/spf13/viper"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var fs embed.FS

// TODO refatorar isso aqui...
func Migrate(ctx context.Context, version int32, dbURL string) error {
	d, err := iofs.New(fs, "sql")
	if err != nil {
		return fmt.Errorf("getting new driver from io/fs: %w", err)
	}

	if dbURL == "" {
		dbURL = fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=%s",
			viper.GetString("VIGIA_DB_USER"),
			viper.GetString("VIGIA_DB_PASSWORD"),
			viper.GetString("VIGIA_DB_HOST"),
			viper.GetString("VIGIA_DB_NAME"),
			viper.GetString("VIGIA_DB_SSLMODE"),
		)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return fmt.Errorf("construct postgres driver: %w", err)
	}

	if version == 0 {
		err = m.Up()
		if err != nil {
			return fmt.Errorf("applying migrations: %w", err)
		}
	} else {
		err = m.Migrate(uint(version))
		if err != nil {
			return fmt.Errorf("applying migrations: %w", err)
		}
	}

	return nil
}
