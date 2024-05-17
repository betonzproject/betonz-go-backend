-- name: GetBets :many
SELECT
	b.*,
	u.username,
	u.role,
	u."etgUsername"
FROM
	"Bet" b
	JOIN "User" u ON u."etgUsername" = b."etgUsername"
WHERE
	(
		sqlc.narg('search')::TEXT IS NULL
		OR b.id::TEXT ILIKE '%' || @search || '%'
		OR u.username ILIKE '%' || @search || '%'
		OR b."providerUsername" ILIKE '%' || @search || '%'
	)
	AND (
		b."productCode" = $1
		OR $1 = 0
	)
	AND (
		b."productType" = $2
		OR $2 = 0
	)
	AND b."startTime" >= sqlc.arg('fromDate')
	AND b."startTime" <= sqlc.arg('toDate')
ORDER BY
	b."startTime" DESC;

-- name: GetTopPayout :many
SELECT
	b1.id,
	u.username,
	u."displayName",
	u."profileImage",
	b1.payout
FROM
	"Bet" b1
	INNER JOIN (
		SELECT
			b."etgUsername",
			MAX(b.payout) AS payout
		FROM
			"Bet" b
		GROUP BY
			b."etgUsername"
		HAVING
			MAX(b.payout) > 0
	) b2 ON b1."etgUsername" = b2."etgUsername" AND b1.payout = b2.payout
	INNER JOIN "User" u ON b1."etgUsername" = u."etgUsername"
WHERE
	(
		@productType::int IS NULL
		OR b1."productType" = @productType
		OR @productType = 0
	)
	AND u.status = 'NORMAL' AND b1."etgUsername" != 'g2rmlmeqvk13' ;

-- name: GetTurnoverByUserId :many
SELECT
	b."productCode",
	sum(b.turnover) AS turnover
FROM
	"Bet" b
	JOIN "User" u USING ("etgUsername")
WHERE
	u.id = $1
	AND b."startTime" >= sqlc.arg('fromDate')
	AND b."startTime" <= sqlc.arg('toDate')
GROUP BY
	b."productCode";

-- name: GetTotalWinLoss :one
SELECT
	COALESCE(sum("winLoss"), 0)::bigint
FROM
	"Bet" b
	JOIN "User" u USING ("etgUsername")
WHERE
	u."role" = 'PLAYER'
	AND b."startTime" >= sqlc.arg('fromDate')
	AND b."startTime" <= sqlc.arg('toDate');

-- name: GetTotalBetAmount :one
SELECT
    COALESCE(sum("bet"), 0)::bigint
FROM
    "Bet" b
	JOIN "User" u USING ("etgUsername")
WHERE
    u.id = $1;

-- name: UpsertBet :exec
INSERT INTO
	"Bet" (
		id,
		"refId",
		"etgUsername",
		"providerUsername",
		"productCode",
		"productType",
		"gameId",
		details,
		turnover,
		bet,
		payout,
		status,
		"startTime",
		"matchTime",
		"endTime",
		"settleTime",
		"progShare",
		"progWin",
		commission,
		"winLoss"
	)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
ON CONFLICT (id) DO
UPDATE
SET
	"refId" = EXCLUDED."refId",
	"etgUsername" = EXCLUDED."etgUsername",
	"providerUsername" = EXCLUDED."providerUsername",
	"productCode" = EXCLUDED."productCode",
	"productType" = EXCLUDED."productType",
	"gameId" = EXCLUDED."gameId",
	details = EXCLUDED.details,
	turnover = EXCLUDED.turnover,
	bet = EXCLUDED.bet,
	payout = EXCLUDED.payout,
	status = EXCLUDED.status,
	"startTime" = EXCLUDED."startTime",
	"matchTime" = EXCLUDED."matchTime",
	"endTime" = EXCLUDED."endTime",
	"settleTime" = EXCLUDED."settleTime",
	"progShare" = EXCLUDED."progShare",
	"progWin" = EXCLUDED."progWin",
	commission = EXCLUDED.commission,
	"winLoss" = EXCLUDED."winLoss";