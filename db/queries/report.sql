-- name: GetDailyPerformance :many
WITH
	TransactionSummary AS (
		-- Get daily sum of approved deposit amount, deposit count, sum of approved withdraw amount, withdrawal count,
		-- sum of bonuses, and sum of withdrawBankFees
		SELECT
			DATE_TRUNC('day', t."updatedAt" AT TIME ZONE 'Asia/Yangon')::DATE AS "date",
			COALESCE(SUM(CASE WHEN t.type = 'DEPOSIT' THEN t.amount ELSE 0 END), 0)::BIGINT AS "depositAmount",
			COALESCE(COUNT(CASE WHEN t.type = 'DEPOSIT' THEN t.* END), 0)::BIGINT AS "depositCount",
			COALESCE(SUM(CASE WHEN t.type = 'WITHDRAW' THEN t.amount ELSE 0 END), 0)::BIGINT AS "withdrawAmount",
			COALESCE(COUNT(CASE WHEN t.type = 'WITHDRAW' THEN t.* END), 0)::BIGINT AS "withdrawCount",
			COALESCE(SUM(t.amount), 0)::BIGINT AS total,
			COALESCE(COUNT(t.*), 0)::BIGINT AS COUNT,
			COALESCE(SUM(t.bonus), 0)::BIGINT AS "bonusGiven",
			COALESCE(SUM(t."withdrawBankFees"), 0)::BIGINT AS "withdrawBankFees"
		FROM
			"TransactionRequest" t
		WHERE
			t."updatedAt" >= sqlc.arg('fromDate')
			AND t."updatedAt" <= sqlc.arg('toDate')
			AND t.status = 'APPROVED'
		GROUP BY
			DATE_TRUNC('day', t."updatedAt" AT TIME ZONE 'Asia/Yangon')
	),
	BetSummary AS (
		-- Get daily sum of win/loss
		SELECT
			DATE_TRUNC('day', "endTime" AT TIME ZONE 'Asia/Yangon')::DATE AS "date",
			COALESCE(SUM("winLoss"), 0)::BIGINT AS "winLoss"
		FROM
			"Bet"
		WHERE
			"endTime" >= sqlc.arg('fromDate')
			AND "endTime" <= sqlc.arg('toDate')
		GROUP BY
			DATE_TRUNC('day', "endTime" AT TIME ZONE 'Asia/Yangon')
	),
	ActivePlayerCount AS (
		-- Get daily active player count
		SELECT
			COALESCE(DATE_TRUNC('day', e."updatedAt" AT TIME ZONE 'Asia/Yangon')) AS "date",
			COUNT(DISTINCT e."userId") AS "activePlayerCount"
		FROM
			"User" u
			JOIN "Event" e ON u.id = e."userId"
			AND e.type = 'ACTIVE'
		WHERE
			u.role = 'PLAYER'
			AND e."updatedAt" >= sqlc.arg('fromDate')
			AND e."updatedAt" <= sqlc.arg('toDate')
		GROUP BY
			COALESCE(DATE_TRUNC('day', e."updatedAt" AT TIME ZONE 'Asia/Yangon'))
	)
SELECT
	COALESCE(bs."date", ts."date", apc."date") AS "createdAt",
	COALESCE(ts."depositAmount", 0),
	COALESCE(ts."depositCount", 0),
	COALESCE(ts."withdrawAmount", 0),
	COALESCE(ts."withdrawCount", 0),
	COALESCE(ts."bonusGiven", 0),
	COALESCE(ts."withdrawBankFees", 0),
	COALESCE(ts."depositAmount", 0),
	COALESCE(bs."winLoss", 0),
	COALESCE(apc."activePlayerCount", 0)
FROM
	TransactionSummary ts
	FULL JOIN BetSummary bs ON ts."date" = bs."date"
	FULL JOIN ActivePlayerCount apc ON ts."date" = apc."date"
ORDER BY
	"createdAt" DESC;
