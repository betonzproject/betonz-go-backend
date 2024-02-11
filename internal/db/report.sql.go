// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: report.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getDailyPerformance = `-- name: GetDailyPerformance :many
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
			t."updatedAt" >= $1
			AND t."updatedAt" <= $2
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
			"endTime" >= $1
			AND "endTime" <= $2
		GROUP BY
			DATE_TRUNC('day', "endTime" AT TIME ZONE 'Asia/Yangon')
	),
	ActivePlayerCount AS (
		-- Get daily active player count
		SELECT
			COALESCE(DATE_TRUNC('day', e."createdAt" AT TIME ZONE 'Asia/Yangon')) AS "date",
			COUNT(*) AS "activePlayerCount"
		FROM
			"User" u
			JOIN "Event" e ON u.id = e."userId"
			AND e.type = 'ACTIVE'
		WHERE
			u.role = 'PLAYER'
			AND e."createdAt" >= $1
			AND e."createdAt" <= $2
		GROUP BY
			COALESCE(DATE_TRUNC('day', e."createdAt" AT TIME ZONE 'Asia/Yangon'))
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
	"createdAt" DESC
`

type GetDailyPerformanceParams struct {
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

type GetDailyPerformanceRow struct {
	CreatedAt         pgtype.Date `json:"createdAt"`
	DepositAmount     int64       `json:"depositAmount"`
	DepositCount      int64       `json:"depositCount"`
	WithdrawAmount    int64       `json:"withdrawAmount"`
	WithdrawCount     int64       `json:"withdrawCount"`
	BonusGiven        int64       `json:"bonusGiven"`
	WithdrawBankFees  int64       `json:"withdrawBankFees"`
	DepositAmount_2   int64       `json:"depositAmount_2"`
	WinLoss           int64       `json:"winLoss"`
	ActivePlayerCount int64       `json:"activePlayerCount"`
}

func (q *Queries) GetDailyPerformance(ctx context.Context, arg GetDailyPerformanceParams) ([]GetDailyPerformanceRow, error) {
	rows, err := q.db.Query(ctx, getDailyPerformance, arg.FromDate, arg.ToDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDailyPerformanceRow
	for rows.Next() {
		var i GetDailyPerformanceRow
		if err := rows.Scan(
			&i.CreatedAt,
			&i.DepositAmount,
			&i.DepositCount,
			&i.WithdrawAmount,
			&i.WithdrawCount,
			&i.BonusGiven,
			&i.WithdrawBankFees,
			&i.DepositAmount_2,
			&i.WinLoss,
			&i.ActivePlayerCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
