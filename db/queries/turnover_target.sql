-- name: GetTurnoverTargetsByUserId :many
WITH "to" AS (
	-- Calculate sum of turnover starting since the start of turnover target
	SELECT
		tt.id,
		sum(b.turnover)::numeric(32, 2) AS "turnoverSoFar"
	FROM
		"TurnoverTarget" tt
		JOIN "TransactionRequest" tr ON tt."transactionRequestId" = tr.id
		JOIN "User" u ON tr."userId" = u.id
		JOIN "Bet" b ON b."etgUsername" = u."etgUsername" AND tr."depositToWallet" = b."productCode"
	WHERE
		b."startTime" >= tt."createdAt"
		AND u.id = $1
	GROUP BY
		tt.id
)
SELECT
	tt.*,
	tr."depositToWallet" AS "productCode",
	COALESCE("turnoverSoFar", 0)::numeric(32, 0) AS "turnoverSoFar"
FROM
	"TurnoverTarget" tt
	LEFT JOIN "to" ON "to".id = tt.id
	JOIN "TransactionRequest" tr ON tt."transactionRequestId" = tr.id
WHERE
	tr."userId" = $1
	AND ("turnoverSoFar" IS NULL OR "turnoverSoFar" < target);

-- name: HasActivePromotionByUserId :one
SELECT EXISTS (
	SELECT
		tt.*
	FROM
		"TurnoverTarget" tt
		JOIN "TransactionRequest" tr ON tt."transactionRequestId" = tr.id
	WHERE
		tr."userId" = $1
		AND tr.promotion = @promotion
);

-- name: HasTurnoverTargetByProductAndUserId :one
SELECT EXISTS (
	SELECT
		tt.*
	FROM
		"TurnoverTarget" tt
		JOIN "TransactionRequest" tr ON tt."transactionRequestId" = tr.id
	WHERE
		tr."userId" = $1
		AND tr."depositToWallet" = sqlc.narg('productCode')
);

-- name: CreateTurnoverTarget :exec
INSERT INTO "TurnoverTarget" (target, "transactionRequestId", "updatedAt") VALUES ($1, $2, now());

-- name: DeleteFulfilledTurnoverTargetsByUserId :exec
DELETE FROM "TurnoverTarget" WHERE id IN (
	-- Get all ids of turnover targets that have been fulfilled
	SELECT
		tt.id
	FROM
		"TurnoverTarget" tt
		JOIN "TransactionRequest" tr ON tt."transactionRequestId" = tr.id
		JOIN "User" u ON tr."userId" = u.id
		JOIN "Bet" b ON b."etgUsername" = u."etgUsername" AND tr."depositToWallet" = b."productCode"
	WHERE
		b."startTime" >= tt."createdAt"
		AND u.id = $1
	GROUP BY
		tt.id
	HAVING
		sum(b.turnover) >= target
);
