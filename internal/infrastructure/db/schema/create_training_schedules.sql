CREATE TABLE IF NOT EXISTS training_schedules (
  id SERIAL PRIMARY KEY,
  trainer_username VARCHAR(100) REFERENCES users(username) ON DELETE CASCADE,
  customer_username VARCHAR(100) REFERENCES customers(username) ON DELETE SET NULL,
  session_id INT REFERENCES customer_sessions(id) ON DELETE SET NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  schedule_type ENUM('APPOINTMENT', 'DAY_OFF') NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);