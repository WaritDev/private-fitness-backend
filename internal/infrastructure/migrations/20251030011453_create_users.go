package migrations
import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createUsers)
}

var createUsers = &Migration{
	Title: "20251030011453_create_users.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(100) PRIMARY KEY,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			role ENUM('ADMIN', 'TRAINER', 'SALES', 'CUSTOMER', 'MANAGER') NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			gender ENUM('MALE', 'FEMALE', 'OTHER') NOT NULL,
			date_of_birth DATE NOT NULL,
			phone_number VARCHAR(20) NOT NULL,
			gmail VARCHAR(255) UNIQUE NOT NULL,
			specialty VARCHAR(100) NULL,
			is_active TINYINT(1) DEFAULT 1,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
		DROP TABLE IF EXISTS users;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
