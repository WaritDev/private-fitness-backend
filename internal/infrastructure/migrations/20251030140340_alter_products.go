package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, alterProducts)
}

var alterProducts = &Migration{
	Title: "20251030140340_alter_products.go",
	Up: func(db *sql.DB) error {
		// Insert default payment account first
		_, err := db.ExecContext(context.Background(), `
		INSERT INTO payment_accounts (account_name, account_number, bank_name, qr_code_image_url, is_active) 
		VALUES ('Default Account', '0000000000', 'Default Bank', 'https://example.com/default-qr.png', 1);
		`)
		if err != nil {
			return err
		}

		// Add payment_account_id column
		_, err = db.ExecContext(context.Background(), `
		ALTER TABLE products 
		ADD COLUMN payment_account_id INT NOT NULL DEFAULT 1 AFTER is_active,
		ADD CONSTRAINT fk_product_payment_account 
		FOREIGN KEY (payment_account_id) REFERENCES payment_accounts(id) ON DELETE CASCADE;
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		// Drop foreign key first
		_, err := db.ExecContext(context.Background(), `
		ALTER TABLE products 
		DROP FOREIGN KEY fk_product_payment_account,
		DROP COLUMN payment_account_id;
		`)
		if err != nil {
			return err
		}

		// Delete default payment account
		_, err = db.ExecContext(context.Background(), `
		DELETE FROM payment_accounts WHERE id = 1;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
