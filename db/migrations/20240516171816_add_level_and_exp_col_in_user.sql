-- migrate:up
ALTER TABLE "User"
ADD COLUMN "level" INT DEFAULT 1 CHECK ("level" <= LEAST(80, "level")),
ADD COLUMN "exp" numeric(65,30) NOT NULL DEFAULT 0;

-- migrate:down
ALTER TABLE "User"
DROP COLUMN "level",
DROP COLUMN "exp";