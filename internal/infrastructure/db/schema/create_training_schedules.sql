CREATE TABLE IF NOT EXISTS `training_schedules` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `trainer_username` VARCHAR(100),
  `customer_username` VARCHAR(100),
  `session_id` INT,
  `start_time` TIMESTAMP NOT NULL,
  `end_time` TIMESTAMP NOT NULL,
  `schedule_type` ENUM('APPOINTMENT','DAY_OFF') NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_sched_trainer` FOREIGN KEY (`trainer_username`) REFERENCES `users`(`username`) ON DELETE CASCADE,
  CONSTRAINT `fk_sched_customer` FOREIGN KEY (`customer_username`) REFERENCES `customers`(`username`) ON DELETE SET NULL,
  CONSTRAINT `fk_sched_session` FOREIGN KEY (`session_id`) REFERENCES `customer_sessions`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;