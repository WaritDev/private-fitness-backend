CREATE TABLE IF NOT EXISTS customer_durations (
  id SERIAL PRIMARY KEY,
  customer_username VARCHAR(100) REFERENCES customers(username) ON DELETE CASCADE,
  sales_username VARCHAR(100) REFERENCES users(username) ON DELETE SET NULL,
  product_id INT REFERENCES products(id) ON DELETE CASCADE,
  purchase_date TIMESTAMP NOT NULL,
  start_date TIMESTAMP NOT NULL,
  end_date TIMESTAMP NOT NULL,
  price_paid DECIMAL(10, 2) NOT NULL,
  discount_amount DECIMAL(10, 2) DEFAULT 0,
  status ENUM('ACTIVE', 'EXPIRED', 'CANCELLED') NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);