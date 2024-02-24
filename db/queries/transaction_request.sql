-- name: GetTransactionRequests :many
SELECT
	tr.*,
	u.username,
	u.role,
	u2.username AS "modifiedByUsername",
	u2.role AS "modifiedByRole"
FROM
	"TransactionRequest" tr
	JOIN "User" u ON u.id = tr."userId"
	LEFT JOIN "User" u2 ON u2.id = tr."modifiedById"
WHERE
	(
		@types::"TransactionType"[] IS NULL
		OR tr."type" = ANY (@types)
	)
	AND (
		@statuses::"TransactionStatus"[] IS NULL
		OR tr.status = ANY (@statuses)
	)
	AND (
		sqlc.narg('search')::TEXT IS NULL
		OR u.username ILIKE '%' || @search || '%'
		OR u2.username ILIKE '%' || @search || '%'
		OR tr.id::TEXT ILIKE '%' || @search || '%'
		OR tr."bankAccountName" ILIKE '%' || @search || '%'
		OR tr."beneficiaryBankAccountName" ILIKE '%' || @search || '%'
		OR tr.remarks ILIKE '%' || @search || '%'
	)
	AND tr."createdAt" >= sqlc.arg('fromDate')
	AND tr."createdAt" <= sqlc.arg('toDate')
ORDER BY
	tr.id DESC;

-- name: GetTransactionRequestById :one
SELECT * FROM "TransactionRequest" WHERE id = $1;

-- name: GetTransactionRequestsByUserId :many
SELECT
	tr.*,
	u.username,
	u.role,
	u2.username AS "modifiedByUsername",
	u2.role AS "modifiedByRole"
FROM
	"TransactionRequest" tr
	JOIN "User" u ON u.id = tr."userId"
	LEFT JOIN "User" u2 ON u2.id = tr."modifiedById"
WHERE
	(
		@types::"TransactionType"[] IS NULL
		OR tr."type" = ANY (@types)
	)
	AND (
		@statuses::"TransactionStatus"[] IS NULL
		OR tr.status = ANY (@statuses)
	)
	AND tr."createdAt" >= sqlc.arg('fromDate')
	AND tr."createdAt" <= sqlc.arg('toDate')
	AND tr."userId" = $1
ORDER BY
	tr.id DESC;

-- name: GetPendingTransactionRequestCount :one
SELECT COUNT(*) FROM "TransactionRequest" WHERE status = 'PENDING';

-- name: HasRecentDepositRequestsByUserId :one
SELECT EXISTS (
	SELECT
		*
	FROM
		"TransactionRequest"
	WHERE
		"userId" = $1
		AND type = 'DEPOSIT'
		AND status = 'PENDING'
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
		AND type = 'WITHDRAW'
		AND status = 'PENDING'
		AND "createdAt" >= now() - INTERVAL '5 minutes'
);

-- name: HasApprovedDepositRequestsWithin30DaysByUserId :one
SELECT EXISTS (
	SELECT
		*
	FROM
		"TransactionRequest"
	WHERE
		"userId" = $1
		AND type = 'DEPOSIT'
		AND status = 'APPROVED'
		AND "updatedAt" >= now() - INTERVAL '30 days'
);

-- name: HasPendingTransactionRequestsWithPromotion :one
SELECT EXISTS (
	SELECT
		*
	FROM
		"TransactionRequest"
	WHERE
		"userId" = $1
		AND type = 'DEPOSIT'
		AND status = 'PENDING'
		AND promotion = $2
);

-- name: GetTotalTransactionAmountAndCount :one
SELECT
	COALESCE(sum(tr.amount), 0)::bigint AS total,
	COALESCE(count(*), 0)::bigint AS count,
	COALESCE(sum(tr.bonus), 0)::bigint AS "bonusTotal"
FROM
	"TransactionRequest" tr
	JOIN "User" u ON tr."userId" = u.id
WHERE
	u."role" = 'PLAYER'
	AND tr."type" = $1
	AND tr.status = 'APPROVED'
	AND tr."updatedAt" >= sqlc.arg('fromDate')
	AND tr."updatedAt" <= sqlc.arg('toDate');

-- name: GetPlayerWithTransactionsCount :one
SELECT
	count(*)
FROM
	(
		SELECT
			DISTINCT "userId"
		FROM
			"TransactionRequest"
		WHERE
			status = 'APPROVED'
			AND "updatedAt" >= sqlc.arg('fromDate')
			AND "updatedAt" <= sqlc.arg('toDate')
	) q;

-- name: GetNewPlayerWithTransactionsCount :one
SELECT
	count(*)
FROM
	(
		SELECT
			DISTINCT "userId"
		FROM
			"TransactionRequest" tr
			JOIN "User" u ON tr."userId" = u.id
		WHERE
			tr.status = 'APPROVED'
			AND tr."updatedAt" >= sqlc.arg('fromDate')
			AND tr."updatedAt" <= sqlc.arg('toDate')
			AND u."createdAt" >= sqlc.arg('fromDate')
			AND u."createdAd" <= sqlc.arg('toDate')
	) q;

-- name: GetBonusRemaining :one
SELECT
	GREATEST(sqlc.arg('limit') - COALESCE(sum(bonus), 0), 0)::numeric(32, 2) AS remaining
FROM
	"TransactionRequest" tr
WHERE
	"userId" = $1
	AND type = 'DEPOSIT'
	AND (status = 'PENDING' OR status = 'APPROVED')
	AND promotion = $2
	AND "updatedAt" >= sqlc.arg('fromDate')
	AND "updatedAt" <= sqlc.arg('toDate');

-- name: CreateTransactionRequest :exec
INSERT INTO
	"TransactionRequest" (
		"userId",
		"bankName",
		"bankAccountName",
		"bankAccountNumber",
		"beneficiaryBankAccountName",
		"beneficiaryBankAccountNumber",
		amount,
		bonus,
		type,
		"depositToWallet",
		promotion,
		"receiptPath",
		status,
		"modifiedById",
		"remarks",
		"updatedAt"
	)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, now());

-- name: UpdateTransactionRequest :exec
UPDATE "TransactionRequest"
SET
	"modifiedById" = $2,
	"receiptPath" = COALESCE($3, "receiptPath"),
	status = $4,
	"withdrawBankFees" = COALESCE($5, 0),
	remarks = $6,
	"updatedAt" = now()
WHERE id = $1;
