CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type ENUM('DURATION', 'SESSION') NOT NULL,
  category ENUM("ECONOMIC", "BUSINESS", "FIRST_CLASS") NOT NULL,
  list_price decimal(10, 2) NOT NULL,
  duration_days INT NULL,
  session_amount INT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);