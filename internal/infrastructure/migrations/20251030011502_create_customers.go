package migrations
import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createCustomers)
}

var createCustomers = &Migration{
	Title: "20251030011502_create_customers.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS customers (
			username VARCHAR(100) PRIMARY KEY,
			health_info TEXT NOT NULL,
			address TEXT NOT NULL,
			company_name VARCHAR(200) NOT NULL,
			company_position VARCHAR(100) NOT NULL,
			marital_status ENUM('SINGLE', 'MARRIED', 'DIVORCED', 'WIDOWED') NOT NULL,
			emergency_contact_name VARCHAR(255) NOT NULL,
			emergency_contact_relationship VARCHAR(50) NOT NULL,
			emergency_contact_phone VARCHAR(20) NOT NULL,
			marketing_source VARCHAR(100) NOT NULL,
			CONSTRAINT fk_customer_user FOREIGN KEY (username)
				REFERENCES users(username) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		DROP TABLE IF EXISTS customers;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
