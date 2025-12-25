-- Migration: Create transfers table
-- Description: Create table for wallet and paylater transfers between users

CREATE TABLE IF NOT EXISTS transfers (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),         -- sender/reseller
  target_user_id BIGINT NOT NULL REFERENCES users(id),  -- receiver
  amount BIGINT NOT NULL CHECK (amount > 0),
  source VARCHAR(20) NOT NULL CHECK (source IN ('wallet', 'paylater')),
  status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'cancelled')),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_transfers_user_id ON transfers(user_id);
CREATE INDEX IF NOT EXISTS idx_transfers_target_user_id ON transfers(target_user_id);
CREATE INDEX IF NOT EXISTS idx_transfers_status ON transfers(status);
CREATE INDEX IF NOT EXISTS idx_transfers_created_at ON transfers(created_at);

-- Add comment
COMMENT ON TABLE transfers IS 'Records of transfers between users from wallet or paylater';
COMMENT ON COLUMN transfers.user_id IS 'ID of the sender/reseller';
COMMENT ON COLUMN transfers.target_user_id IS 'ID of the receiver';
COMMENT ON COLUMN transfers.amount IS 'Transfer amount in smallest currency unit';
COMMENT ON COLUMN transfers.source IS 'Source of funds: wallet or paylater';
COMMENT ON COLUMN transfers.status IS 'Transfer status: pending, completed, failed, cancelled';
