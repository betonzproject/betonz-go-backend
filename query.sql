-- name: GetUserById :one
SELECT id, username, role, email, "displayName", "phoneNumber", "mainWallet", dob, "profileImage", "isEmailVerified" FROM "User" WHERE id = $1;

-- name: GetExtendedUserByUsername :one
SELECT * FROM "User" WHERE username = $1 AND (@roles::"Role"[] IS NULL OR role = ANY(@roles));



-- name: GetTransactionRequests :many
SELECT
	tr.*,
	u.username,
	u.role,
	u2.username AS "modifiedByUsername",
	u2.role AS "modifiedByRole"
FROM
	"TransactionRequest" tr
JOIN "User" u ON
	u.id = tr."userId"
LEFT JOIN "User" u2 ON
	u2.id = tr."modifiedById"
WHERE 
	(@types::"TransactionType"[] IS NULL OR tr."type" = ANY(@types))
AND
	(@statuses::"TransactionStatus"[] IS NULL OR tr.status = ANY(@statuses))
AND (
	sqlc.narg('search')::text IS NULL
	OR u.username ILIKE '%' || @search || '%' 
	OR u2.username ILIKE '%' || @search || '%'
	OR tr."bankAccountName" ILIKE '%' || @search || '%'
	OR tr."beneficiaryBankAccountName" ILIKE '%' || @search || '%'
	OR tr.remarks ILIKE '%' || @search || '%'
)
AND
	tr."createdAt" >= sqlc.arg('fromDate')
AND
	tr."createdAt" < sqlc.arg('toDate')
ORDER BY
	tr.id DESC;
