package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createPaymentVerifications)
}

var createPaymentVerifications = &Migration{
	Title: "20251031000000_create_payment_verifications.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS payment_verifications (
			id INT AUTO_INCREMENT PRIMARY KEY,
			customer_username VARCHAR(100) NOT NULL,
			product_id INT NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			slip_file_path VARCHAR(500),
			slip_id VARCHAR(100),
			verification_status ENUM('PENDING','VERIFIED','REJECTED') NOT NULL DEFAULT 'PENDING',
			slip2go_response TEXT,
			verified_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_payment_customer FOREIGN KEY (customer_username) REFERENCES users(username) ON DELETE CASCADE,
			CONSTRAINT fk_payment_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
			INDEX idx_customer_username (customer_username),
			INDEX idx_verification_status (verification_status),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		DROP TABLE IF EXISTS payment_verifications;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
