-- migrate:up
ALTER TABLE "Bank"
ADD CONSTRAINT unique_account_bank UNIQUE ("accountNumber", name);

-- migrate:down
ALTER TABLE "Bank"
DROP CONSTRAINT IF EXISTS unique_account_bank;
