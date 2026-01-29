CREATE TABLE discount_requests (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES users(id),
    discount_percentage DECIMAL(5,2) NOT NULL,
    reason TEXT NOT NULL,
    status discount_status NOT NULL,
    rule_id BIGINT,
    approved_by_id BIGINT REFERENCES users(id),
    approval_comment TEXT DEFAULT 'Not Updated by manager',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
