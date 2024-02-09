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
	AND u.status = 'NORMAL';

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
