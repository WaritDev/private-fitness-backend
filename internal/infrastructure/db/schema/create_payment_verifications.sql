CREATE TABLE IF NOT EXISTS `payment_verifications` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `customer_username` VARCHAR(100) NOT NULL,
  `product_id` INT NOT NULL,
  `amount` DECIMAL(10,2) NOT NULL,
  `slip_file_path` VARCHAR(500),
  `slip_id` VARCHAR(100),
  `verification_status` ENUM('PENDING','VERIFIED','REJECTED') NOT NULL DEFAULT 'PENDING',
  `slip2go_response` TEXT,
  `verified_at` TIMESTAMP NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_payment_customer` FOREIGN KEY (`customer_username`) REFERENCES `users`(`username`) ON DELETE CASCADE,
  CONSTRAINT `fk_payment_product` FOREIGN KEY (`product_id`) REFERENCES `products`(`id`) ON DELETE CASCADE,
  INDEX `idx_customer_username` (`customer_username`),
  INDEX `idx_verification_status` (`verification_status`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
