package migrations

import (
	"context"
	"database/sql"
	"fmt"
)

func init() {
	Migrations = append(Migrations, alterUsersUsernameCaseSensitive)
}

var alterUsersUsernameCaseSensitive = &Migration{
	Title: "20251102094911_alter_users_username_case_sensitive.go",
	Up: func(db *sql.DB) error {
		// Step 1: Drop ALL foreign key constraints that reference username columns
		dropStatements := []string{
			// customers -> users
			`ALTER TABLE customers DROP FOREIGN KEY fk_customer_user`,
			// customer_sessions -> customers, users
			`ALTER TABLE customer_sessions DROP FOREIGN KEY fk_sess_customer`,
			`ALTER TABLE customer_sessions DROP FOREIGN KEY fk_sess_trainer`,
			`ALTER TABLE customer_sessions DROP FOREIGN KEY fk_sess_sales`,
			// customer_durations -> customers, users
			`ALTER TABLE customer_durations DROP FOREIGN KEY fk_dur_customer`,
			`ALTER TABLE customer_durations DROP FOREIGN KEY fk_dur_sales`,
			// customer_logs -> customers
			`ALTER TABLE customer_logs DROP FOREIGN KEY fk_log_customer`,
			// training_schedules -> users, customers
			`ALTER TABLE training_schedules DROP FOREIGN KEY fk_sched_trainer`,
			`ALTER TABLE training_schedules DROP FOREIGN KEY fk_sched_customer`,
			`ALTER TABLE training_schedules DROP FOREIGN KEY fk_sched_session`,
			// training_availabilities -> users
			`ALTER TABLE training_availabilities DROP FOREIGN KEY fk_avail_trainer`,
		}
		for _, stmt := range dropStatements {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				fmt.Printf("Warning: %v\n", err)
			}
		}

		// Step 2: Change ALL username columns collation to be case-sensitive
		alterStatements := []string{
			// users table (primary)
			`ALTER TABLE users CHANGE COLUMN username username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL`,
			// customers table
			`ALTER TABLE customers CHANGE COLUMN username username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL`,
			// customer_logs table
			`ALTER TABLE customer_logs CHANGE COLUMN customer_username customer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL`,
			// customer_sessions table
			`ALTER TABLE customer_sessions CHANGE COLUMN customer_username customer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL`,
			`ALTER TABLE customer_sessions CHANGE COLUMN trainer_username trainer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL`,
			`ALTER TABLE customer_sessions CHANGE COLUMN sales_username sales_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL`,
			// customer_durations table
			`ALTER TABLE customer_durations CHANGE COLUMN customer_username customer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL`,
			`ALTER TABLE customer_durations CHANGE COLUMN sales_username sales_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL`,
			// training_schedules table
			`ALTER TABLE training_schedules CHANGE COLUMN trainer_username trainer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL`,
			`ALTER TABLE training_schedules CHANGE COLUMN customer_username customer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL`,
			// training_availabilities table
			`ALTER TABLE training_availabilities CHANGE COLUMN trainer_username trainer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL`,
		}

		for _, stmt := range alterStatements {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				fmt.Printf("Warning (alter): %v\n", err)
			}
		}

		// Step 3: Re-add foreign key constraints
		addStatements := []string{
			// customers -> users
			`ALTER TABLE customers ADD CONSTRAINT fk_customer_user FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE`,
			// customer_sessions -> customers, users
			`ALTER TABLE customer_sessions ADD CONSTRAINT fk_sess_customer FOREIGN KEY (customer_username) REFERENCES customers(username) ON DELETE CASCADE`,
			`ALTER TABLE customer_sessions ADD CONSTRAINT fk_sess_trainer FOREIGN KEY (trainer_username) REFERENCES users(username) ON DELETE SET NULL`,
			`ALTER TABLE customer_sessions ADD CONSTRAINT fk_sess_sales FOREIGN KEY (sales_username) REFERENCES users(username) ON DELETE SET NULL`,
			// customer_durations -> customers, users
			`ALTER TABLE customer_durations ADD CONSTRAINT fk_dur_customer FOREIGN KEY (customer_username) REFERENCES customers(username) ON DELETE CASCADE`,
			`ALTER TABLE customer_durations ADD CONSTRAINT fk_dur_sales FOREIGN KEY (sales_username) REFERENCES users(username) ON DELETE SET NULL`,
			// customer_logs -> customers
			`ALTER TABLE customer_logs ADD CONSTRAINT fk_log_customer FOREIGN KEY (customer_username) REFERENCES customers(username) ON DELETE CASCADE`,
			// training_schedules -> users, customers
			`ALTER TABLE training_schedules ADD CONSTRAINT fk_sched_trainer FOREIGN KEY (trainer_username) REFERENCES users(username) ON DELETE CASCADE`,
			`ALTER TABLE training_schedules ADD CONSTRAINT fk_sched_customer FOREIGN KEY (customer_username) REFERENCES customers(username) ON DELETE SET NULL`,
			// training_availabilities -> users
			`ALTER TABLE training_availabilities ADD CONSTRAINT fk_avail_trainer FOREIGN KEY (trainer_username) REFERENCES users(username) ON DELETE CASCADE`,
		}

		for _, stmt := range addStatements {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				fmt.Printf("Warning (add FK): %v\n", err)
			}
		}

		return nil
	},
	Down: func(db *sql.DB) error {
		// Step 1: Drop ALL foreign key constraints
		dropStatements := []string{
			// customers -> users
			`ALTER TABLE customers DROP FOREIGN KEY fk_customer_user`,
			// customer_sessions -> customers, users
			`ALTER TABLE customer_sessions DROP FOREIGN KEY fk_sess_customer`,
			`ALTER TABLE customer_sessions DROP FOREIGN KEY fk_sess_trainer`,
			`ALTER TABLE customer_sessions DROP FOREIGN KEY fk_sess_sales`,
			// customer_durations -> customers, users
			`ALTER TABLE customer_durations DROP FOREIGN KEY fk_dur_customer`,
			`ALTER TABLE customer_durations DROP FOREIGN KEY fk_dur_sales`,
			// customer_logs -> customers
			`ALTER TABLE customer_logs DROP FOREIGN KEY fk_log_customer`,
			// training_schedules -> users, customers
			`ALTER TABLE training_schedules DROP FOREIGN KEY fk_sched_trainer`,
			`ALTER TABLE training_schedules DROP FOREIGN KEY fk_sched_customer`,
			`ALTER TABLE training_schedules DROP FOREIGN KEY fk_sched_session`,
			// training_availabilities -> users
			`ALTER TABLE training_availabilities DROP FOREIGN KEY fk_avail_trainer`,
		}
		for _, stmt := range dropStatements {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				fmt.Printf("Warning: %v\n", err)
			}
		}

		// Step 2: Change back ALL username columns to case-insensitive
		alterStatements := []string{
			`ALTER TABLE users CHANGE COLUMN username username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL`,
			`ALTER TABLE customers CHANGE COLUMN username username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL`,
			`ALTER TABLE customer_logs CHANGE COLUMN customer_username customer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL`,
			`ALTER TABLE customer_sessions CHANGE COLUMN customer_username customer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL`,
			`ALTER TABLE customer_sessions CHANGE COLUMN trainer_username trainer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL`,
			`ALTER TABLE customer_sessions CHANGE COLUMN sales_username sales_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL`,
			`ALTER TABLE customer_durations CHANGE COLUMN customer_username customer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL`,
			`ALTER TABLE customer_durations CHANGE COLUMN sales_username sales_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL`,
			`ALTER TABLE training_schedules CHANGE COLUMN trainer_username trainer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL`,
			`ALTER TABLE training_schedules CHANGE COLUMN customer_username customer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL`,
			`ALTER TABLE training_availabilities CHANGE COLUMN trainer_username trainer_username VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL`,
		}

		for _, stmt := range alterStatements {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				fmt.Printf("Warning (alter): %v\n", err)
			}
		}

		// Step 3: Re-add foreign key constraints
		addStatements := []string{
			// customers -> users
			`ALTER TABLE customers ADD CONSTRAINT fk_customer_user FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE`,
			// customer_sessions -> customers, users
			`ALTER TABLE customer_sessions ADD CONSTRAINT fk_sess_customer FOREIGN KEY (customer_username) REFERENCES customers(username) ON DELETE CASCADE`,
			`ALTER TABLE customer_sessions ADD CONSTRAINT fk_sess_trainer FOREIGN KEY (trainer_username) REFERENCES users(username) ON DELETE SET NULL`,
			`ALTER TABLE customer_sessions ADD CONSTRAINT fk_sess_sales FOREIGN KEY (sales_username) REFERENCES users(username) ON DELETE SET NULL`,
			// customer_durations -> customers, users
			`ALTER TABLE customer_durations ADD CONSTRAINT fk_dur_customer FOREIGN KEY (customer_username) REFERENCES customers(username) ON DELETE CASCADE`,
			`ALTER TABLE customer_durations ADD CONSTRAINT fk_dur_sales FOREIGN KEY (sales_username) REFERENCES users(username) ON DELETE SET NULL`,
			// customer_logs -> customers
			`ALTER TABLE customer_logs ADD CONSTRAINT fk_log_customer FOREIGN KEY (customer_username) REFERENCES customers(username) ON DELETE CASCADE`,
			// training_schedules -> users, customers
			`ALTER TABLE training_schedules ADD CONSTRAINT fk_sched_trainer FOREIGN KEY (trainer_username) REFERENCES users(username) ON DELETE CASCADE`,
			`ALTER TABLE training_schedules ADD CONSTRAINT fk_sched_customer FOREIGN KEY (customer_username) REFERENCES customers(username) ON DELETE SET NULL`,
			// training_availabilities -> users
			`ALTER TABLE training_availabilities ADD CONSTRAINT fk_avail_trainer FOREIGN KEY (trainer_username) REFERENCES users(username) ON DELETE CASCADE`,
		}

		for _, stmt := range addStatements {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				fmt.Printf("Warning (add FK): %v\n", err)
			}
		}

		return nil
	},
}
