CREATE TABLE IF NOT EXISTS `customer_logs` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `customer_username` VARCHAR(100),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `status` ENUM('PENDING','CONFIRMED') NULL,
  `schedule_id` INT NULL,
  `log_type` ENUM('CHECK_IN','CHECK_OUT','BOOK_SESSION','CANCEL_SESSION') NOT NULL,
  CONSTRAINT `fk_log_schedule` FOREIGN KEY (`schedule_id`) REFERENCES `training_schedules`(`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_log_customer` FOREIGN KEY (`customer_username`) REFERENCES `customers`(`username`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;