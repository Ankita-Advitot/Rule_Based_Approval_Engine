CREATE TYPE expense_status AS ENUM (
  'APPLIED',
  'AUTO_APPROVED',
  'PENDING',
  'APPROVED',
  'REJECTED',
  'AUTO_REJECTED',
  'CANCELLED'
);

CREATE TABLE expense_requests (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    category VARCHAR(100) NOT NULL,
    reason TEXT NOT NULL,
    status expense_status NOT NULL,
    rule_id BIGINT,
    approved_by_id BIGINT REFERENCES users(id),
    approval_comment TEXT DEFAULT 'Not Updated by manager',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
