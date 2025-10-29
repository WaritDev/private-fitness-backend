CREATE TABLE IF NOT EXISTS training_availabilities (
  id SERIAL PRIMARY KEY,
  trainer_username VARCHAR(100) REFERENCES users(username) ON DELETE CASCADE,
  day_of_week ENUM('MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY') NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL
);