package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, dropEmailFromUsers)
}

var dropEmailFromUsers = &Migration{
	Title: "20251030184825_drop_email_from_users.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
			ALTER TABLE users DROP COLUMN email;
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
			ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL UNIQUE AFTER password;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
