-- name: GetVerificationTokenByHash :one
SELECT
	vt.*,
	u.username,
	u.email
FROM
	"VerificationToken" vt
	LEFT JOIN "User" u ON vt."userId" = u.id
WHERE
	"tokenHash" = $1;

-- name: UpsertVerificationToken :exec
INSERT INTO
	"VerificationToken" ("tokenHash", "userId", "registerInfo", "updatedAt")
VALUES
	($1, $2, $3, now())
ON CONFLICT ("userId") DO
UPDATE
SET
	"tokenHash" = excluded."tokenHash",
	"registerInfo" = excluded."registerInfo",
	"createdAt" = now(),
	"updatedAt" = now();

-- name: DeleteVerificationTokenByHash :exec
DELETE FROM "VerificationToken" WHERE "tokenHash" = $1;
