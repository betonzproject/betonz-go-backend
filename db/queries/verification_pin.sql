-- name: GetVerificationPinByUserId :one
SELECT
	prt.*
FROM
	"VerificationPin" prt
	JOIN "User" u ON prt."userId" = u.id
WHERE
	"userId" = $1;

-- name: GetVerificationPinByPin :one
SELECT
	prt.*,
	u.username,
	u.email
FROM
	"VerificationPin" prt
	JOIN "User" u ON prt."userId" = u.id
WHERE
	"pin" = $1;

-- name: UpsertVerificationPin :exec
INSERT INTO
	"VerificationPin" AS prt ("pin", "userId", "updatedAt")
VALUES
	($1, $2, now())
ON CONFLICT ("userId") DO
UPDATE
SET
	"pin" = excluded."pin",
	"createdAt" = now(),
	"updatedAt" = now();

-- name: DeleteVerificationPin :exec
DELETE FROM "VerificationPin"
WHERE
	"pin" = $1;
