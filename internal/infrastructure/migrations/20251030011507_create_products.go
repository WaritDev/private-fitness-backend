package migrations
import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createProducts)
}

var createProducts = &Migration{
	Title: "20251030011507_create_products.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS products (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name TEXT NOT NULL,
			type ENUM('DURATION', 'SESSION') NOT NULL,
			category ENUM('ECONOMIC', 'BUSINESS', 'FIRST_CLASS') NOT NULL,
			list_price DECIMAL(10,2) NOT NULL,
			duration_days INT NULL,
			session_amount INT NULL,
			is_active TINYINT(1) DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		DROP TABLE IF EXISTS products;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
