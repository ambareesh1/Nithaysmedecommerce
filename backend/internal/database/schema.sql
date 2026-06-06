CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    order_number VARCHAR(100) NOT NULL,
    total_amount NUMERIC(10, 2) NOT NULL,
    total_quantity INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    subtotal NUMERIC(10, 2) NOT NULL
);

INSERT INTO products (name, description, category, price, stock, image_url) VALUES
('Digital Thermometer', 'A digital thermometer used to measure body temperature.', 'Diagnostics', 499, 25, 'https://example.com/thermometer.png'),
('Surgical Gloves', 'Disposable latex surgical gloves for medical use.', 'Protection', 299, 100, 'https://example.com/gloves.png'),
('Blood Pressure Monitor', 'Automatic blood pressure monitor for home use.', 'Diagnostics', 1999, 15, 'https://example.com/bp.png'),
('Pulse Oximeter', 'Fingertip pulse oximeter to measure oxygen levels.', 'Diagnostics', 899, 40, 'https://example.com/oximeter.png'),
('Nebulizer Machine', 'Compact nebulizer machine for respiratory care.', 'Equipment', 2499, 10, 'https://example.com/nebulizer.png'),
('First Aid Kit', 'Complete first aid kit for emergency care.', 'Supplies', 799, 50, 'https://example.com/firstaid.png')
ON CONFLICT DO NOTHING;
