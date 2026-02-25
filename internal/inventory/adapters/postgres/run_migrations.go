package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations executes the migrations against the provided DB connection
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {

	poolConn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer poolConn.Release()

	rawConn := poolConn.Conn()

	migrator, err := migrate.NewMigrator(ctx, rawConn, "schema_version")
	if err != nil {
		return fmt.Errorf("unable to create migrator: %w", err)
	}

	content, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	if err := migrator.LoadMigrations(content); err != nil {
		return fmt.Errorf("unable to load migrations: %w", err)
	}

	fmt.Printf("DEBUG: Loaded %d migrations.\n", len(migrator.Migrations))

	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("unable to migrate: %w", err)
	}

	ver, _ := migrator.GetCurrentVersion(ctx)
	fmt.Printf("Migration successful. Current version: %d\n", ver)

	return nil
}
