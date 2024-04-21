package db

import "database/sql"

type dbMigration interface {
	Up(*sql.DB) error
}

type createPagesTableMigration struct{}

func (m *createPagesTableMigration) Up(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS pages (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			link TEXT NOT NULL,
			created_at INTEGER NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}
