package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, alterUsers)
}

var alterUsers = &Migration{
	Title: "20251030184131_alter_users.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		ALTER TABLE users DROP COLUMN email;
		`)
		return err
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		ALTER TABLE users 
		ADD COLUMN email VARCHAR(255) UNIQUE NOT NULL AFTER phone_number;
		`)
		return err
	},
}