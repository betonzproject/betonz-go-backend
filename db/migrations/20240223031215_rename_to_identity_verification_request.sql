-- migrate:up
ALTER TABLE "IdentityVerificationRequests" ALTER COLUMN dob SET NOT NULL;
ALTER TABLE "IdentityVerificationRequests" RENAME TO "IdentityVerificationRequest";
ALTER TYPE "IdentityVerificationStatus" ADD value 'INCOMPLETE';

-- migrate:down
CREATE TYPE "IdentityVerificationStatus_new" AS ENUM('VERIFIED', 'REJECTED', 'PENDING');
ALTER TABLE "IdentityVerificationRequest" ALTER COLUMN status TYPE "IdentityVerificationStatus_new" USING (status::TEXT::"IdentityVerificationStatus_new");
ALTER TYPE "IdentityVerificationStatus" RENAME TO "IdentityVerificationStatus_old";
ALTER TYPE "IdentityVerificationStatus_new" RENAME TO "IdentityVerificationStatus";
DROP TYPE "IdentityVerificationStatus_old";

ALTER TABLE "IdentityVerificationRequest" ALTER COLUMN dob DROP NOT NULL;
ALTER TABLE "IdentityVerificationRequest" RENAME TO "IdentityVerificationRequests";
