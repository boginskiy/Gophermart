-- Обратное изменение таблицы orders
ALTER TABLE orders ALTER COLUMN accrual TYPE INTEGER USING accrual::INTEGER;

-- Обратное изменение таблицы loyalty_orders
ALTER TABLE loyalty_orders ALTER COLUMN deduction TYPE INTEGER USING deduction::INTEGER;