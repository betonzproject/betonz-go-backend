-- name: GetBankById :one
SELECT * FROM "Bank" WHERE id = $1;

-- name: GetBanksByUserId :many
SELECT * FROM "Bank" WHERE "userId" = $1 ORDER BY "createdAt", "accountName";

-- name: GetSystemBankById :one
SELECT b.* FROM "Bank" b JOIN "User" u ON b."userId" = u.id WHERE u.role = 'SYSTEM' AND b.id = $1;

-- name: GetSystemBanks :many
SELECT b.* FROM "Bank" b JOIN "User" u ON b."userId" = u.id WHERE u.role = 'SYSTEM' ORDER BY b."createdAt", b."accountName";

-- name: GetSystemBanksByBankName :many
SELECT b.* FROM "Bank" b JOIN "User" u ON b."userId" = u.id WHERE u.role = 'SYSTEM' AND b.name = $1;

-- name: CreateBank :one
INSERT INTO
	"Bank" (id, "userId", name, "accountName", "accountNumber", "updatedAt")
VALUES
	(gen_random_uuid(), $1, $2, $3, $4, now())
RETURNING
	*;

-- name: CreateSystemBank :exec
INSERT INTO
	"Bank" (id, "userId", name, "accountName", "accountNumber", "updatedAt")
SELECT
	gen_random_uuid(),
	id,
	$1,
	$2,
	$3,
	now()
FROM
	"User"
WHERE
	role = 'SYSTEM'
LIMIT
	1;

-- name: UpdateBank :exec
UPDATE "Bank"
SET
	name = $2,
	"accountName" = COALESCE(sqlc.narg('accountName'), "accountName"),
	"accountNumber" = COALESCE(sqlc.narg('accountNumber'), "accountNumber"),
	"updatedAt" = now()
WHERE
	id = $1;

-- name: UpdateSystemBank :exec
UPDATE "Bank"
SET
	"accountName" = COALESCE(sqlc.narg('accountName'), "accountName"),
	"accountNumber" = COALESCE(sqlc.narg('accountNumber'), "accountNumber"),
	disabled = @disabled,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: DeleteBankById :one
DELETE FROM "Bank" b USING "User" u WHERE b."userId" = u.id AND b.id = $1 RETURNING b.*;

-- name: DeleteSystemBankById :exec
DELETE FROM "Bank" b USING "User" u WHERE b."userId" = u.id AND b.id = $1 AND u.role = 'SYSTEM';
