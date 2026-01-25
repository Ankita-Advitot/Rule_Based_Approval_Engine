INSERT INTO rules (request_type, condition, action, grade_id)
VALUES
('LEAVE', '{"max_days": 5}', 'AUTO_APPROVE', 2);

INSERT INTO rules (request_type, condition, action, grade_id)
VALUES
('LEAVE', '{"max_days": 10}', 'AUTO_APPROVE', 3);

INSERT INTO rules (request_type, condition, action, grade_id)
VALUES
('EXPENSE', '{"max_amount": 10000}', 'AUTO_APPROVE', 2);

INSERT INTO rules (request_type, condition, action, grade_id)
VALUES
('EXPENSE', '{"max_amount": 20000}', 'AUTO_APPROVE', 3);

INSERT INTO rules (request_type, condition, action, grade_id)
VALUES
('DISCOUNT', '{"max_percent": 20}', 'AUTO_APPROVE', 2);

INSERT INTO rules (request_type, condition, action, grade_id)
VALUES
('DISCOUNT', '{"max_percent": 30}', 'AUTO_APPROVE', 3);