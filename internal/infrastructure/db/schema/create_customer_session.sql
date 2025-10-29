CREATE TABLE IF NOT EXISTS customer_sessions (
  id SERIAL PRIMARY KEY,
  customer_username VARCHAR(100) REFERENCES customers(username) ON DELETE CASCADE,
  trainer_username VARCHAR(100) REFERENCES users(username) ON DELETE SET NULL,
  sales_username VARCHAR(100) REFERENCES users(username) ON DELETE SET NULL,
  product_id INT REFERENCES products(id) ON DELETE CASCADE,
  purchase_date TIMESTAMP NOT NULL,
  total_sessions INT NOT NULL,
  used_sessions INT DEFAULT 0,
  price_paid DECIMAL(10, 2) NOT NULL,
  discount_amount DECIMAL(10, 2) DEFAULT 0,
  status ENUM('ACTIVE', 'EXPIRED', 'CANCELLED', 'COMPLETED') NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);