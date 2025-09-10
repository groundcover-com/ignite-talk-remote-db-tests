-- Create users table migration
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert some test data
INSERT INTO users (email, name) VALUES 
    ('john@example.com', 'John Doe'),
    ('jane@example.com', 'Jane Smith')
ON CONFLICT (email) DO NOTHING;

