package models

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type Service struct {
	PassService
	db     *sql.DB
	ctx    context.Context // background context
	cancel func()          // cancel background context
	// Datasource name.
	DSN string
}

func NewServices(dsn string) (*Service, error) {
	service := &Service{
		DSN: dsn,
	}

	service.ctx, service.cancel = context.WithCancel(context.Background())

	err := service.Open()
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	err = service.migrate()
	if err != nil {
		return nil, fmt.Errorf("could not migrate files: %w", err)
	}

	service.PassService = NewPassService(service.db)

	return service, nil
}

// The code to migrate sql files was take from Ben Johnson's WTF repo: https://github.com/benbjohnson/wtf/blob/main/sqlite/sqlite.go

// Open opens the database connection.
func (s *Service) Open() (err error) {
	// Ensure a DSN is set before attempting to open the database.
	if s.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	// Make the parent directory unless using an in-memory db.
	if s.DSN != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(s.DSN), 0700); err != nil {
			return err
		}
	}

	// Connect to the database.
	if s.db, err = sql.Open("sqlite3", s.DSN); err != nil {
		return err
	}

	// Enable WAL. SQLite performs better with the WAL  because it allows
	// multiple readers to operate while data is being written.
	if _, err := s.db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	// Enable foreign key checks. For historical reasons, SQLite does not check
	// foreign key constraints by default... which is kinda insane. There's some
	// overhead on inserts to verify foreign key integrity but it's definitely
	// worth it.
	if _, err := s.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	if err := s.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

// migrate sets up migration tracking and executes pending migration files.
//
// Migration files are embedded in the sqlite/migration folder and are executed
// in lexigraphical order.
//
// Once a migration is run, its name is stored in the 'migrations' table so it
// is not re-executed. Migrations run in a transaction to prevent partial
// migrations.
func (s *Service) migrate() error {
	// Ensure the 'migrations' table exists so we don't duplicate migrations.
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	// Read migration files from our embedded file system.
	// This uses Go 1.16's 'embed' package.
	names, err := fs.Glob(migrationFS, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	// Loop over all migration files and execute them in order.
	for _, name := range names {
		if err := s.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}
	return nil
}

// migrate runs a single migration file within a transaction. On success, the
// migration file name is saved to the "migrations" table to prevent re-running.
func (s *Service) migrateFile(name string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure migration has not already been run.
	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // already run migration, skip
	}

	// Read and execute migration file.
	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	// Insert record into migrations to prevent re-running migration.
	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES (?)`, name); err != nil {
		return err
	}

	return tx.Commit()
}

// Close closes the database connection.
func (s *Service) Close() error {
	// Cancel background context.
	s.cancel()

	// Close database.
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
