package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createPaymentAccounts)
}

var createPaymentAccounts = &Migration{
	Title: "20251030140331_create_payment_accounts.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS payment_accounts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			account_name VARCHAR(255) NOT NULL,
			account_number VARCHAR(50) NOT NULL,
			bank_name VARCHAR(100) NOT NULL,
			qr_code_image_url TEXT NOT NULL,
			is_active TINYINT(1) NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		DROP TABLE IF EXISTS payment_accounts;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
