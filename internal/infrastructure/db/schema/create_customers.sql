CREATE TABLE IF NOT EXISTS customers (
    username VARCHAR(100) PRIMARY KEY REFERENCES users(username) ON DELETE CASCADE,
    health_info TEXT NOT NULL,
    address TEXT NOT NULL,
    company_name VARCHAR(200) NOT NULL,
    company_position VARCHAR(100) NOT NULL,
    marital_status TEXT CHECK (marital_status IN ('SINGLE', 'MARRIED', 'DIVORCED', 'WIDOWED')) NOT NULL,
    emergency_contact_name VARCHAR(255) NOT NULL,
    emergency_contact_relationship VARCHAR(50) NOT NULL,
    emergency_contact_phone VARCHAR(20) NOT NULL,
    marketing_source VARCHAR(100) NOT NULL
);