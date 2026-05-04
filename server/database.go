package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func openDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := prepareUsersTable(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func prepareUsersTable(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			password_hash TEXT,
			deleted_at TEXT
		)
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	if err := addDeletedAtColumnIfMissing(db); err != nil {
		return err
	}
	if err := addPasswordHashColumnIfMissing(db); err != nil {
		return err
	}
	if err := createActiveUserNameIndex(db); err != nil {
		return err
	}

	return nil
}

func addDeletedAtColumnIfMissing(db *sql.DB) error {
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		return fmt.Errorf("read users table info: %w", err)
	}
	defer rows.Close()

	var hasDeletedAt bool
	for rows.Next() {
		var cid int
		var name string
		var columnType string
		var notNull int
		var defaultValue any
		var pk int

		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &pk); err != nil {
			return fmt.Errorf("scan users table info: %w", err)
		}
		if name == "deleted_at" {
			hasDeletedAt = true
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate users table info: %w", err)
	}

	if hasDeletedAt {
		return nil
	}

	if _, err := db.Exec("ALTER TABLE users ADD COLUMN deleted_at TEXT"); err != nil {
		return fmt.Errorf("add deleted_at column: %w", err)
	}

	return nil
}

func addPasswordHashColumnIfMissing(db *sql.DB) error {
	hasPasswordHash, err := hasUsersColumn(db, "password_hash")
	if err != nil {
		return err
	}
	if hasPasswordHash {
		return nil
	}

	if _, err := db.Exec("ALTER TABLE users ADD COLUMN password_hash TEXT"); err != nil {
		return fmt.Errorf("add password_hash column: %w", err)
	}

	return nil
}

func createActiveUserNameIndex(db *sql.DB) error {
	const query = `
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name_active
		ON users(name)
		WHERE deleted_at IS NULL
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("create active user name index: %w", err)
	}

	return nil
}

func hasUsersColumn(db *sql.DB, columnName string) (bool, error) {
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		return false, fmt.Errorf("read users table info: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var columnType string
		var notNull int
		var defaultValue any
		var pk int

		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &pk); err != nil {
			return false, fmt.Errorf("scan users table info: %w", err)
		}
		if name == columnName {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate users table info: %w", err)
	}

	return false, nil
}
