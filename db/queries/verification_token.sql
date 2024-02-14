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

-- name: CreateVerificationToken :exec
INSERT INTO "VerificationToken" ("tokenHash", "registerInfo", "updatedAt") VALUES ($1, $2, now());

-- name: DeleteVerificationTokenByHash :exec
DELETE FROM "VerificationToken" WHERE "tokenHash" = $1;
