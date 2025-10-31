CREATE TABLE IF NOT EXISTS `training_availabilities` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `trainer_username` VARCHAR(100) NOT NULL,
  `day_of_week` ENUM('MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY') NOT NULL,
  `start_time` TIMESTAMP NOT NULL,
  `end_time` TIMESTAMP NOT NULL,
  CONSTRAINT `fk_avail_trainer` FOREIGN KEY (`trainer_username`) REFERENCES `users`(`username`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;