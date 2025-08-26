-- Create categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for better performance
CREATE INDEX idx_categories_name ON categories(name);

-- Insert default categories
INSERT INTO categories (name) VALUES 
('Food & Beverages'),
('Electronics'),
('Clothing'),
('Books'),
('Health & Beauty'),
('Others');