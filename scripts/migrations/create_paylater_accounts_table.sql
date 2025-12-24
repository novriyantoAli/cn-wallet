-- Migration: Create paylater_accounts table
-- Description: This table stores paylater account information for users
-- Created: 2025-12-24

CREATE TABLE paylater_accounts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
  credit_limit BIGINT NOT NULL CHECK (credit_limit > 0),
  outstanding BIGINT NOT NULL DEFAULT 0 CHECK (outstanding >= 0),
  available_limit BIGINT NOT NULL,
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CHECK (available_limit = credit_limit - outstanding)
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_paylater_accounts_user_id ON paylater_accounts(user_id);

-- Create index on status for filtering
CREATE INDEX idx_paylater_accounts_status ON paylater_accounts(status);

-- Comments
COMMENT ON TABLE paylater_accounts IS 'Stores paylater account information for users';
COMMENT ON COLUMN paylater_accounts.id IS 'Primary key';
COMMENT ON COLUMN paylater_accounts.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN paylater_accounts.credit_limit IS 'Maximum credit limit for the account';
COMMENT ON COLUMN paylater_accounts.outstanding IS 'Current outstanding balance';
COMMENT ON COLUMN paylater_accounts.available_limit IS 'Available credit (credit_limit - outstanding)';
COMMENT ON COLUMN paylater_accounts.status IS 'Account status: active, suspended, or closed';
COMMENT ON COLUMN paylater_accounts.created_at IS 'Timestamp when account was created';
