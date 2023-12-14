-- name: GetTransactionLogs :many
SELECT
	tr.*,
	u.username AS "initiatorUsername",
	u.role AS "initiatorRole",
	u2.username AS "beneficiaryUsername",
	u2.role AS "beneficiaryRole"
FROM
	"Transaction" tr
JOIN "User" u ON
	u.id = tr."initiatorId"
JOIN "User" u2 ON
	u2.id = tr."beneficiaryId"
WHERE 
(
	sqlc.narg('search')::text IS NULL
	OR u.username ILIKE '%' || @search || '%' 
	OR u2.username ILIKE '%' || @search || '%'
	OR tr.remarks ILIKE '%' || @search || '%'
)
AND
	tr."createdAt" >= sqlc.arg('fromDate')
AND
	tr."createdAt" < sqlc.arg('toDate')
ORDER BY
	tr.id DESC;