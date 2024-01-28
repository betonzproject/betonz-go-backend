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
	tr."createdAt" <= sqlc.arg('toDate')
ORDER BY
	tr.id DESC;

-- name: GetTransactionRequestsByUserId :many
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
AND
	tr."createdAt" >= sqlc.arg('fromDate')
AND
	tr."createdAt" <= sqlc.arg('toDate')
AND
	tr."userId" = $1
ORDER BY
	tr.id DESC;

-- name: HasRecentDepositRequestsByUserId :one
SELECT EXISTS (
	SELECT
		*
	FROM
		"TransactionRequest"
	WHERE
		"userId" = $1
	AND type = 'DEPOSIT'::"TransactionType"
		AND status = 'PENDING'::"TransactionStatus"
		AND "createdAt" >= now() - INTERVAL '1 minute'
);

-- name: HasRecentWithdrawRequestsByUserId :one
SELECT EXISTS (
	SELECT
		*
	FROM
		"TransactionRequest"
	WHERE
		"userId" = $1
	AND type = 'WITHDRAW'::"TransactionType"
		AND status = 'PENDING'::"TransactionStatus"
		AND "createdAt" >= now() - INTERVAL '5 minutes'
);



-- name: CreateTransactionRequest :exec
INSERT INTO "TransactionRequest" (
	"userId",
	"bankName",
	"bankAccountName",
	"bankAccountNumber",
	"beneficiaryBankAccountName",
	"beneficiaryBankAccountNumber",
	amount,
	bonus,
	type,
	"receiptPath",
	status,
	"updatedAt"
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now());
