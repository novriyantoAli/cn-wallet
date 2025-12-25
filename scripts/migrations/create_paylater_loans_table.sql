-- Migration: Create paylater_loans table
-- Description: This table stores loan records for users with paylater accounts
-- Created: 2025-12-25

CREATE TABLE paylater_loans (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  amount BIGINT NOT NULL CHECK (amount > 0),
  interest BIGINT NOT NULL CHECK (interest >= 0),
  total BIGINT NOT NULL CHECK (total = amount + interest),
  due_date DATE NOT NULL,
  source VARCHAR(20) NOT NULL CHECK (source IN ('transfer', 'checkout')),
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_paylater_loans_user_id ON paylater_loans(user_id);

-- Create index on status for filtering
CREATE INDEX idx_paylater_loans_status ON paylater_loans(status);

-- Create index on due_date for querying overdue loans
CREATE INDEX idx_paylater_loans_due_date ON paylater_loans(due_date);

-- Create composite index for user_id and status
CREATE INDEX idx_paylater_loans_user_status ON paylater_loans(user_id, status);

-- Comments
COMMENT ON TABLE paylater_loans IS 'Stores loan records for users with paylater accounts';
COMMENT ON COLUMN paylater_loans.id IS 'Primary key';
COMMENT ON COLUMN paylater_loans.user_id IS 'Foreign key to users table (reseller)';
COMMENT ON COLUMN paylater_loans.amount IS 'Original loan amount';
COMMENT ON COLUMN paylater_loans.interest IS 'Interest charged on the loan';
COMMENT ON COLUMN paylater_loans.total IS 'Total amount to be repaid (amount + interest)';
COMMENT ON COLUMN paylater_loans.due_date IS 'Date when the loan is due';
COMMENT ON COLUMN paylater_loans.source IS 'Source of the loan: transfer or checkout';
COMMENT ON COLUMN paylater_loans.status IS 'Loan status: pending, active, paid, overdue, or defaulted';
COMMENT ON COLUMN paylater_loans.created_at IS 'Timestamp when loan was created';
