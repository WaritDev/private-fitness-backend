CREATE TABLE IF NOT EXISTS `customer_durations` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `customer_username` VARCHAR(100),
  `sales_username` VARCHAR(100),
  `product_id` INT,
  `purchase_date` TIMESTAMP NOT NULL,
  `start_date` TIMESTAMP NOT NULL,
  `end_date` TIMESTAMP NOT NULL,
  `price_paid` DECIMAL(10,2) NOT NULL,
  `discount_amount` DECIMAL(10,2) DEFAULT 0,
  `status` ENUM('ACTIVE','EXPIRED','CANCELLED') NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_dur_customer` FOREIGN KEY (`customer_username`) REFERENCES `customers`(`username`) ON DELETE CASCADE,
  CONSTRAINT `fk_dur_sales` FOREIGN KEY (`sales_username`) REFERENCES `users`(`username`) ON DELETE SET NULL,
  CONSTRAINT `fk_dur_product` FOREIGN KEY (`product_id`) REFERENCES `products`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;