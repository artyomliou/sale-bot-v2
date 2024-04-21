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

type addNotificationTitleMigration struct{}

func (m *addNotificationTitleMigration) Up(db *sql.DB) error {
	var exists int
	row := db.QueryRow(`select count(*) from pragma_table_info('pages') where name = 'notification_title';`)
	if err := row.Scan(&exists); err != nil {
		return err
	}
	if exists == 1 {
		return nil
	}

	_, err := db.Exec(`
		ALTER TABLE pages ADD notification_title TEXT;
	`)
	if err != nil {
		return err
	}
	return nil
}
