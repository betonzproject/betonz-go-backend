-- migrate:up
ALTER TABLE "TransactionRequest"
ALTER COLUMN "bankName" DROP NOT NULL,
ALTER COLUMN "bankAccountName" DROP NOT NULL,
ALTER COLUMN "bankAccountNumber" DROP NOT NULL;

-- migrate:down
ALTER TABLE "TransactionRequest"
ALTER COLUMN "bankName" SET NOT NULL,
ALTER COLUMN "bankAccountName" SET NOT NULL,
ALTER COLUMN "bankAccountNumber" SET NOT NULL;
