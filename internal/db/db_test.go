package db_test

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"artyomliou/sale-bot-v2/internal/db"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {
	file, err := os.CreateTemp("", "db.*.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	conn, err := db.NewConnection(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	page := &crawlers.Page{
		ID:    "test",
		Title: "test",
		Link:  "test",
	}

	exists, err := conn.CheckPageExists(page)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, exists)

	err = conn.CreatePage(page)
	if err != nil {
		t.Fatal(err)
	}

	exists, err = conn.CheckPageExists(page)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, exists)
}
