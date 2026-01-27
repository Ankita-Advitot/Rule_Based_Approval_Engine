-- Expense wallet
ALTER TABLE expense
ADD CONSTRAINT expense_user_id_unique UNIQUE (user_id);

-- Leave wallet
ALTER TABLE leaves
ADD CONSTRAINT leaves_user_id_unique UNIQUE (user_id);

-- Discount wallet
ALTER TABLE discount
ADD CONSTRAINT discount_user_id_unique UNIQUE (user_id);
