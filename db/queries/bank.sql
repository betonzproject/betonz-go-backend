-- name: GetSystemBanks :many
SELECT
	b.id, b.name, b."accountName", b."accountNumber", b.disabled
FROM
	"Bank" b
JOIN "User" u ON
	b."userId" = u.id
WHERE
	u.ROLE = 'SYSTEM'::"Role"
ORDER BY
	b."createdAt", b."accountName";



-- name: CreateSystemBank :exec
INSERT INTO "Bank" ("id", "userId", "name", "accountName", "accountNumber", "updatedAt")
SELECT gen_random_uuid(), id, $1, $2, $3, now() FROM "User" WHERE role = 'SYSTEM'::"Role" LIMIT 1;



-- name: UpdateSystemBank :exec
UPDATE
	"Bank"
SET
	"accountName" = COALESCE(sqlc.narg('accountName'), "accountName"),
	"accountNumber" = COALESCE(sqlc.narg('accountNumber'), "accountNumber"),
	disabled = @disabled,
	"updatedAt" = now()
WHERE
	id = $1;



-- name: DeleteSystemBankById :exec
DELETE FROM "Bank" b USING "User" u WHERE b."userId" = u.id AND b.id = $1 AND u.role = 'SYSTEM'::"Role";
