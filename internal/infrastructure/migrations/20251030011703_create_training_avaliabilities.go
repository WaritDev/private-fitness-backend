package migrations
import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createTrainingAvaliabilities)
}

var createTrainingAvaliabilities = &Migration{
	Title: "20251030011703_create_training_avaliabilities.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS training_availabilities (
			id INT AUTO_INCREMENT PRIMARY KEY,
			trainer_username VARCHAR(100) NOT NULL,
			day_of_week ENUM('MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY') NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			CONSTRAINT fk_avail_trainer FOREIGN KEY (trainer_username)
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
		DROP TABLE IF EXISTS training_availabilities;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
