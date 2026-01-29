CREATE TYPE leave_status AS ENUM (
  'APPLIED',
  'AUTO_APPROVED',
  'PENDING',
  'APPROVED',
  'REJECTED',
  'AUTO_REJECTED',
  'CANCELLED'
);

CREATE TABLE leave_requests (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES users(id),
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    reason TEXT NOT NULL,
    leave_type VARCHAR(20) NOT NULL,
    status leave_status NOT NULL,
    approved_by_id BIGINT REFERENCES users(id),
    rule_id BIGINT,
    approval_comment TEXT DEFAULT 'Not Updated by manager',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
