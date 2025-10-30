CREATE TABLE IF NOT EXISTS `payment_accounts` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `account_name` VARCHAR(255) NOT NULL,
  `account_number` VARCHAR(50) NOT NULL,
  `bank_name` VARCHAR(100) NOT NULL,
  `qr_code_image_url` TEXT NOT NULL,
  `is_active` TINYINT(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;