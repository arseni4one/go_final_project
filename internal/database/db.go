package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbFile string) error {
	var err error

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Создаём таблицу (IF NOT EXISTS)
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(255) NOT NULL,
		comment TEXT,
		repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
	`
	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}
