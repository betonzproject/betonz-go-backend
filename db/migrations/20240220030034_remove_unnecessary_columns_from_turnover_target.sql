-- migrate:up
ALTER TABLE "TurnoverTarget"
DROP COLUMN "userId",
DROP COLUMN "productCode",
DROP COLUMN "promoCode";

-- migrate:down
ALTER TABLE "TurnoverTarget"
ADD COLUMN "userId" UUID,
ADD COLUMN "productCode" integer,
ADD COLUMN "promoCode" "PromotionType";

UPDATE "TurnoverTarget"
SET
	"userId" = u.id,
	"productCode" = tr."depositToWallet",
	"promoCode" = tr."promotion"
FROM
	"TurnoverTarget" tt
	JOIN "TransactionRequest" tr ON tt."transactionRequestId" = tr.id
	JOIN "User" u ON tr."userId" = u.id;

ALTER TABLE ONLY "TurnoverTarget"
ADD CONSTRAINT "TurnoverTarget_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User" (id) ON UPDATE CASCADE ON DELETE RESTRICT,
ALTER COLUMN "userId" SET NOT NULL,
ALTER COLUMN "productCode" SET NOT NULL;
