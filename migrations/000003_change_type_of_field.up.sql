-- Изменение таблицы orders
ALTER TABLE orders ALTER COLUMN accrual TYPE FLOAT USING accrual::FLOAT;

-- Изменение таблицы loyalty_orders
ALTER TABLE loyalty_orders ALTER COLUMN deduction TYPE FLOAT USING deduction::FLOAT;