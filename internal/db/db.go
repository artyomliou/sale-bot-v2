package db

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"database/sql"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DbConnection struct {
	db *sql.DB
}

func NewConnection(dbPath string) (*DbConnection, error) {
	dbDir := path.Dir(dbPath)

	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	dbMigrations := []dbMigration{
		&createPagesTableMigration{},
		&addNotificationTitleMigration{},
	}
	for _, m := range dbMigrations {
		if err := m.Up(db); err != nil {
			return nil, err
		}
	}

	return &DbConnection{
		db: db,
	}, nil
}

func (conn *DbConnection) Close() error {
	return conn.db.Close()
}

func (conn *DbConnection) CheckPageExists(page *crawlers.Page) (bool, error) {
	stmt, err := conn.db.Prepare(`SELECT EXISTS(SELECT 1 FROM pages WHERE id=?)`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var exists bool
	if err := stmt.QueryRow(page.ID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (conn *DbConnection) CreatePage(page *crawlers.Page) error {
	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into pages(id, title, link, created_at) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(page.ID, page.Title, page.Link, time.Now().Unix()); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
