-- migrate:up
CREATE TYPE "VipType" AS ENUM (
    'BRONZE',
    'SILVER',
    'GOLD',
    'PLATINUM_I',
    'PLATINUM_II',
    'PLATINUM_III',
    'PLATINUM_IV',
    'DIAMOND_I',
    'DIAMOND_II',
    'DIAMOND_III',
    'DIAMOND_IV',
    'JADE',
    'KYAWTHUITE'
);

ALTER TABLE "User" ADD COLUMN "vipLevel" "VipType" DEFAULT 'BRONZE';

-- migrate:down
-- Remove vipLevel column and drop VipType enum
ALTER TABLE "User" DROP COLUMN "vipLevel";
DROP TYPE "VipType";
