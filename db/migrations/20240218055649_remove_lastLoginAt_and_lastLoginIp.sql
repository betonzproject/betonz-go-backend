-- migrate:up
ALTER TABLE "User"
DROP COLUMN "lastLoginAt",
DROP COLUMN "lastLoginIp";

-- migrate:down
ALTER TABLE "User"
ADD COLUMN "lastLoginIp" TEXT,
ADD COLUMN "lastLoginAt" timestamp(3) WITHOUT TIME ZONE;
