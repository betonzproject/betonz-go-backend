-- migrate:up
ALTER TABLE "TransactionRequest"
ADD COLUMN "transactionNo" text;


-- migrate:down
ALTER TABLE "TransactionRequest"
DROP COLUMN "transactionNo";