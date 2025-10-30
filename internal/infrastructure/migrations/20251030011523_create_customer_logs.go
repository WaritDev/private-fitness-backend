package migrations
import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createCustomerLogs)
}

var createCustomerLogs = &Migration{
	Title: "20251030011523_create_customer_logs.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS customer_logs (
			id INT AUTO_INCREMENT PRIMARY KEY,
			customer_username VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			log_type ENUM('CHECK_IN','CHECK_OUT','BOOK_SESSION','CANCEL_SESSION') NOT NULL,
			CONSTRAINT fk_log_customer FOREIGN KEY (customer_username)
				REFERENCES customers(username)
				ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		DROP TABLE IF EXISTS customer_logs;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
