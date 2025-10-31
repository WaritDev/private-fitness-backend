package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, alterCustomerLogsAddCheckinFields)
}

var alterCustomerLogsAddCheckinFields = &Migration{
	Title: "20251101011128_alter_customer_logs_add_checkin_fields.go",
	Up: func(db *sql.DB) error {
		// Step 1: Add status column for pending/confirmed check-ins
		_, err := db.ExecContext(context.Background(), `
			ALTER TABLE customer_logs 
			ADD COLUMN status ENUM('PENDING', 'CONFIRMED') NULL AFTER log_type
		`)
		if err != nil {
			return err
		}

		// Step 2: Add schedule_id column to link with training_schedules
		_, err = db.ExecContext(context.Background(), `
			ALTER TABLE customer_logs 
			ADD COLUMN schedule_id INT NULL AFTER status,
			ADD CONSTRAINT fk_log_schedule FOREIGN KEY (schedule_id)
				REFERENCES training_schedules(id)
				ON DELETE SET NULL
		`)
		if err != nil {
			return err
		}

		// Step 3: Set default status for existing CHECK_IN logs as CONFIRMED
		_, err = db.ExecContext(context.Background(), `
			UPDATE customer_logs 
			SET status = 'CONFIRMED' 
			WHERE log_type = 'CHECK_IN' AND status IS NULL
		`)
		if err != nil {
			return err
		}

		return nil
	},
	Down: func(db *sql.DB) error {
		// Rollback: Drop columns
		_, err := db.ExecContext(context.Background(), `
			ALTER TABLE customer_logs 
			DROP FOREIGN KEY fk_log_schedule,
			DROP COLUMN schedule_id,
			DROP COLUMN status
		`)
		if err != nil {
			return err
		}

		return nil
	},
}
