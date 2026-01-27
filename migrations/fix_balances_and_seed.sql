BEGIN;
-- LEAVES
INSERT INTO leaves (user_id, total_allocated, remaining_count)
VALUES
  (2, 36, 36),
  (3, 45, 45)
ON CONFLICT (user_id) DO NOTHING;

-- EXPENSE
INSERT INTO expense (user_id, total_amount, remaining_amount)
VALUES
  (2, 25000, 25000),
  (3, 50000, 50000)
ON CONFLICT (user_id) DO NOTHING;

-- DISCOUNT
INSERT INTO discount (user_id, total_discount, remaining_discount)
VALUES
  (2, 30, 30),
  (3 ,30, 30)
ON CONFLICT (user_id) DO NOTHING;

COMMIT;



