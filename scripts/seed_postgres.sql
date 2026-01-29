-- PostgreSQL/Supabase Seed Data
-- GoPOS API Initial Data

-- Admin user (password: Admin123!)
INSERT INTO users (id, email, password_hash, name, phone, role, is_active, created_at, updated_at)
VALUES (
    'a1b2c3d4-e5f6-7890-abcd-ef1234567890',
    'admin@gopos.local',
    '$2a$10$ZvsrDTGNLZ5KiGqJtWLqCug8T7GJPnkIRlHix6ECS0FwpXE2d7..2', -- Admin123!
    'Administrator',
    '081234567890',
    'admin',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (id) DO NOTHING;

-- Manager user (password: Manager123!)
INSERT INTO users (id, email, password_hash, name, phone, role, is_active, created_at, updated_at)
VALUES (
    'b2c3d4e5-f678-90ab-cdef-123456789012',
    'manager@gopos.local',
    '$2a$10$oEgvZnGBHfuDiUN6Vq2zdOISOKNpk6jNOjJTM9bpX5HUbTxldVqiO', -- Manager123!
    'Store Manager',
    '081234567891',
    'manager',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (id) DO NOTHING;

-- Cashier user (password: Cashier123!)
INSERT INTO users (id, email, password_hash, name, phone, role, is_active, created_at, updated_at)
VALUES (
    'c3d4e5f6-7890-abcd-ef12-345678901234',
    'cashier@gopos.local',
    '$2a$10$/tOYhYjQrLhe/eTAhmc7wuHTat/ltvF/6pZEEuLxfhXAZ4d7o70oq', -- Cashier123!
    'Cashier Staff',
    '081234567892',
    'cashier',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (id) DO NOTHING;

-- Categories
INSERT INTO categories (id, name, description, slug, is_active, created_at, updated_at)
VALUES 
    ('cat-11111111-1111-1111-1111-111111111111', 'Beverages', 'Drinks and beverages', 'beverages', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('cat-22222222-2222-2222-2222-222222222222', 'Food', 'Food items and snacks', 'food', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('cat-33333333-3333-3333-3333-333333333333', 'Electronics', 'Electronic devices and accessories', 'electronics', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('cat-44444444-4444-4444-4444-444444444444', 'Household', 'Household items and supplies', 'household', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Products
INSERT INTO products (id, category_id, sku, name, description, price, stock, image_url, is_active, created_at, updated_at)
VALUES 
    -- Beverages
    ('prod-1111-1111-1111-111111111111', 'cat-11111111-1111-1111-1111-111111111111', 'BEV-001', 'Mineral Water 600ml', 'Fresh mineral water', 5000, 100, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('prod-1111-1111-1111-222222222222', 'cat-11111111-1111-1111-1111-111111111111', 'BEV-002', 'Orange Juice 350ml', 'Fresh orange juice', 12000, 50, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('prod-1111-1111-1111-333333333333', 'cat-11111111-1111-1111-1111-111111111111', 'BEV-003', 'Coffee Latte', 'Premium coffee latte', 25000, 30, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    
    -- Food
    ('prod-2222-2222-2222-111111111111', 'cat-22222222-2222-2222-2222-222222222222', 'FOD-001', 'Chocolate Bar', 'Delicious chocolate bar', 15000, 75, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('prod-2222-2222-2222-222222222222', 'cat-22222222-2222-2222-2222-222222222222', 'FOD-002', 'Potato Chips', 'Crispy potato chips', 18000, 60, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('prod-2222-2222-2222-333333333333', 'cat-22222222-2222-2222-2222-222222222222', 'FOD-003', 'Instant Noodles', 'Quick instant noodles', 8000, 120, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    
    -- Electronics
    ('prod-3333-3333-3333-111111111111', 'cat-33333333-3333-3333-3333-333333333333', 'ELC-001', 'USB Cable', 'USB Type-C cable 1m', 35000, 40, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('prod-3333-3333-3333-222222222222', 'cat-33333333-3333-3333-3333-333333333333', 'ELC-002', 'Power Bank 10000mAh', 'Portable power bank', 150000, 25, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    
    -- Household
    ('prod-4444-4444-4444-111111111111', 'cat-44444444-4444-4444-4444-444444444444', 'HH-001', 'Tissue Box', 'Soft facial tissue', 12000, 80, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('prod-4444-4444-4444-222222222222', 'cat-44444444-4444-4444-4444-444444444444', 'HH-002', 'Hand Soap', 'Antibacterial hand soap', 22000, 45, '', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Customers
INSERT INTO customers (id, name, email, phone, address, loyalty_points, created_at, updated_at)
VALUES 
    ('cust-1111-1111-1111-111111111111', 'John Doe', 'john@example.com', '081111222333', 'Jl. Sudirman No. 1', 150, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('cust-2222-2222-2222-222222222222', 'Jane Smith', 'jane@example.com', '081444555666', 'Jl. Gatot Subroto No. 2', 250, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('cust-3333-3333-3333-333333333333', 'Bob Wilson', 'bob@example.com', '081777888999', 'Jl. Thamrin No. 3', 50, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;
