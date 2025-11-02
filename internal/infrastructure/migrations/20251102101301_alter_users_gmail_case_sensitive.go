package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, alterUsersGmailCaseSensitive)
}

var alterUsersGmailCaseSensitive = &Migration{
	Title: "20251102101301_alter_users_gmail_case_sensitive.go",
	Up: func(db *sql.DB) error {
		// Change gmail column collation to be case-sensitive
		_, err := db.ExecContext(context.Background(), `
		ALTER TABLE users 
		CHANGE COLUMN gmail gmail VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		// Rollback: Change back to case-insensitive
		_, err := db.ExecContext(context.Background(), `
		ALTER TABLE users 
		CHANGE COLUMN gmail gmail VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL
		`)
		if err != nil {
			return err
		}
		return nil
	},
}

