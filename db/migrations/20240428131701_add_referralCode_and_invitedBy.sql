-- migrate:up
ALTER TABLE "User"
ADD COLUMN "referralCode" VARCHAR(7),
ADD COLUMN "invitedBy" VARCHAR(7);


-- migrate:down
ALTER TABLE "User"
DROP COLUMN "referralCode",
DROP COLUMN "invitedBy";