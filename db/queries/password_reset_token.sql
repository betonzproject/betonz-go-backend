-- name: GetPasswordResetTokenByUserId :one
SELECT prt.* FROM "PasswordResetToken" prt JOIN "User" u ON prt."userId" = u.id WHERE "userId" = $1;

-- name: GetPasswordResetTokenByHash :one
SELECT prt.*, u.username, u.email FROM "PasswordResetToken" prt JOIN "User" u ON prt."userId" = u.id WHERE "tokenHash" = $1;

-- name: CreatePasswordResetToken :exec
INSERT INTO "PasswordResetToken" ("tokenHash", "userId", "updatedAt") VALUES ($1, $2, now());

-- name: DeletePasswordResetToken :exec
DELETE FROM "PasswordResetToken" WHERE "tokenHash" = $1;
