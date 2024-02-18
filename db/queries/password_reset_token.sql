-- name: GetPasswordResetTokenByUserId :one
SELECT prt.* FROM "PasswordResetToken" prt JOIN "User" u ON prt."userId" = u.id WHERE "userId" = $1;

-- name: GetPasswordResetTokenByHash :one
SELECT prt.*, u.username, u.email FROM "PasswordResetToken" prt JOIN "User" u ON prt."userId" = u.id WHERE "tokenHash" = $1;

-- name: UpsertPasswordResetToken :exec
INSERT INTO
	"PasswordResetToken" AS prt ("tokenHash", "userId", "updatedAt")
VALUES
	($1, $2, now())
ON CONFLICT ("userId") DO
UPDATE
SET
	"tokenHash" = excluded."tokenHash",
	"createdAt" = now(),
	"updatedAt" = now();

-- name: DeletePasswordResetToken :exec
DELETE FROM "PasswordResetToken" WHERE "tokenHash" = $1;
