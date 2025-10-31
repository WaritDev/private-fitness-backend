package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, alterTrainingAvailabilitiesTimeToTimestamp)
}

var alterTrainingAvailabilitiesTimeToTimestamp = &Migration{
	Title: "20251031120000_alter_training_availabilities_time_to_timestamp.go",
	Up: func(db *sql.DB) error {
		// Step 1: Add new TIMESTAMP columns
		_, err := db.ExecContext(context.Background(), `
			ALTER TABLE training_availabilities 
			ADD COLUMN start_time_new TIMESTAMP NULL AFTER start_time,
			ADD COLUMN end_time_new TIMESTAMP NULL AFTER end_time
		`)
		if err != nil {
			return err
		}

		// Step 2: Migrate data from TIME to TIMESTAMP (use today's date + time)
		_, err = db.ExecContext(context.Background(), `
			UPDATE training_availabilities 
			SET start_time_new = TIMESTAMP(CONCAT(CURDATE(), ' ', start_time)),
				end_time_new = TIMESTAMP(CONCAT(CURDATE(), ' ', end_time))
		`)
		if err != nil {
			return err
		}

		// Step 3: Drop old TIME columns
		_, err = db.ExecContext(context.Background(), `
			ALTER TABLE training_availabilities 
			DROP COLUMN start_time,
			DROP COLUMN end_time
		`)
		if err != nil {
			return err
		}

		// Step 4: Rename new columns to original names
		_, err = db.ExecContext(context.Background(), `
			ALTER TABLE training_availabilities 
			CHANGE COLUMN start_time_new start_time TIMESTAMP NOT NULL,
			CHANGE COLUMN end_time_new end_time TIMESTAMP NOT NULL
		`)
		if err != nil {
			return err
		}

		return nil
	},
	Down: func(db *sql.DB) error {
		// Rollback: Convert TIMESTAMP back to TIME
		_, err := db.ExecContext(context.Background(), `
			ALTER TABLE training_availabilities 
			ADD COLUMN start_time_old TIME NULL AFTER start_time,
			ADD COLUMN end_time_old TIME NULL AFTER end_time
		`)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(context.Background(), `
			UPDATE training_availabilities 
			SET start_time_old = TIME(start_time),
				end_time_old = TIME(end_time)
		`)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(context.Background(), `
			ALTER TABLE training_availabilities 
			DROP COLUMN start_time,
			DROP COLUMN end_time
		`)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(context.Background(), `
			ALTER TABLE training_availabilities 
			CHANGE COLUMN start_time_old start_time TIME NOT NULL,
			CHANGE COLUMN end_time_old end_time TIME NOT NULL
		`)
		if err != nil {
			return err
		}

		return nil
	},
}
