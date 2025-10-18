-- Создание таблицы loyalty_orders
CREATE TABLE loyalty_orders (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL CHECK(code ~* '^[0-9]+$'),
    deduction INTEGER,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
);
