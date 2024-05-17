// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: ranking_system.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getUserBenefits = `-- name: GetUserBenefits :one
SELECT
    "vipLevel",
    CASE
        WHEN "vipLevel" = 'PLATINUM_I' THEN 200000
        WHEN "vipLevel" = 'PLATINUM_II' THEN 210000
        WHEN "vipLevel" = 'PLATINUM_III' THEN 220000
        WHEN "vipLevel" = 'PLATINUM_IV' THEN 230000
        WHEN "vipLevel" = 'DIAMOND_I' THEN 300000
        WHEN "vipLevel" = 'DIAMOND_II' THEN 310000
        WHEN "vipLevel" = 'DIAMOND_III' THEN 320000
        WHEN "vipLevel" = 'DIAMOND_IV' THEN 330000
        WHEN "vipLevel" = 'JADE' THEN 500000
        WHEN "vipLevel" = 'KYAWTHUITE' THEN 777000
        ELSE 0
    END AS birthdayBonus,
    CASE
        WHEN "vipLevel" IN ('PLATINUM_I', 'PLATINUM_II', 'PLATINUM_III', 'PLATINUM_IV', 'DIAMOND_I', 'DIAMOND_II', 'DIAMOND_III', 'DIAMOND_IV', 'JADE', 'KYAWTHUITE') THEN 'YES'
        ELSE 'NO'
    END AS birthdayGift,
    CASE
        WHEN "vipLevel" IN ('PLATINUM_I', 'PLATINUM_II', 'PLATINUM_III', 'PLATINUM_IV', 'DIAMOND_I', 'DIAMOND_II', 'DIAMOND_III', 'DIAMOND_IV', 'JADE', 'KYAWTHUITE') THEN 'YES'
        ELSE 'NO'
    END AS monthlyGift,
    CASE
        WHEN "vipLevel" IN ('PLATINUM_I', 'PLATINUM_II', 'PLATINUM_III', 'PLATINUM_IV', 'DIAMOND_I', 'DIAMOND_II', 'DIAMOND_III', 'DIAMOND_IV', 'JADE', 'KYAWTHUITE') THEN 'YES'
        ELSE 'NO'
    END AS anniversaryGift
FROM
    "User"
WHERE
    id = $1
`

type GetUserBenefitsRow struct {
	VipLevel        NullVipType `json:"vipLevel"`
	Birthdaybonus   int32       `json:"birthdaybonus"`
	Birthdaygift    string      `json:"birthdaygift"`
	Monthlygift     string      `json:"monthlygift"`
	Anniversarygift string      `json:"anniversarygift"`
}

func (q *Queries) GetUserBenefits(ctx context.Context, id pgtype.UUID) (GetUserBenefitsRow, error) {
	row := q.db.QueryRow(ctx, getUserBenefits, id)
	var i GetUserBenefitsRow
	err := row.Scan(
		&i.VipLevel,
		&i.Birthdaybonus,
		&i.Birthdaygift,
		&i.Monthlygift,
		&i.Anniversarygift,
	)
	return i, err
}

const getUserRanking = `-- name: GetUserRanking :one
SELECT
    u.id,
    u.username,
    u."vipLevel",
    COALESCE(SUM(b.bet), 0)::bigint AS totalBetAmount,
    CASE
        WHEN COALESCE(SUM(b.bet), 0) >= 140000000 THEN 'KYAWTHUITE'
        WHEN COALESCE(SUM(b.bet), 0) >= 120000000 THEN 'JADE'
        WHEN COALESCE(SUM(b.bet), 0) >= 80000000 THEN 'DIAMOND_IV'
        WHEN COALESCE(SUM(b.bet), 0) >= 70000000 THEN 'DIAMOND_III'
        WHEN COALESCE(SUM(b.bet), 0) >= 60000000 THEN 'DIAMOND_II'
        WHEN COALESCE(SUM(b.bet), 0) >= 50000000 THEN 'DIAMOND_I'
        WHEN COALESCE(SUM(b.bet), 0) >= 40000000 THEN 'PLATINUM_IV'
        WHEN COALESCE(SUM(b.bet), 0) >= 35000000 THEN 'PLATINUM_III'
        WHEN COALESCE(SUM(b.bet), 0) >= 30000000 THEN 'PLATINUM_II'
        WHEN COALESCE(SUM(b.bet), 0) >= 25000000 THEN 'PLATINUM_I'
        WHEN COALESCE(SUM(b.bet), 0) >= 15000000 THEN 'GOLD'
        WHEN COALESCE(SUM(b.bet), 0) >= 7000000 THEN 'SILVER'
        ELSE 'BRONZE'
    END AS newVipLevel,
    CASE
        WHEN COALESCE(SUM(b.bet), 0) >= 140000000 THEN 100::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 120000000 THEN ((COALESCE(SUM(b.bet), 0) - 120000000) * 100 / (140000000 - 120000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 80000000 THEN ((COALESCE(SUM(b.bet), 0) - 80000000) * 100 / (120000000 - 80000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 70000000 THEN ((COALESCE(SUM(b.bet), 0) - 70000000) * 100 / (80000000 - 70000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 60000000 THEN ((COALESCE(SUM(b.bet), 0) - 60000000) * 100 / (70000000 - 60000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 50000000 THEN ((COALESCE(SUM(b.bet), 0) - 50000000) * 100 / (60000000 - 50000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 40000000 THEN ((COALESCE(SUM(b.bet), 0) - 40000000) * 100 / (50000000 - 40000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 35000000 THEN ((COALESCE(SUM(b.bet), 0) - 35000000) * 100 / (40000000 - 35000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 30000000 THEN ((COALESCE(SUM(b.bet), 0) - 30000000) * 100 / (35000000 - 30000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 25000000 THEN ((COALESCE(SUM(b.bet), 0) - 25000000) * 100 / (30000000 - 25000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 15000000 THEN ((COALESCE(SUM(b.bet), 0) - 15000000) * 100 / (25000000 - 15000000))::numeric(5, 2)
        WHEN COALESCE(SUM(b.bet), 0) >= 7000000 THEN ((COALESCE(SUM(b.bet), 0) - 7000000) * 100 / (15000000 - 7000000))::numeric(5, 2)
        ELSE (COALESCE(SUM(b.bet), 0) * 100 / 7000000)::numeric(5, 2)
    END AS rankProgress
FROM
    "User" u
LEFT JOIN
    "Bet" b ON u."etgUsername" = b."etgUsername"
WHERE
    u.id = $1
GROUP BY
    u.id, u.username, u."vipLevel"
`

type GetUserRankingRow struct {
	ID             pgtype.UUID    `json:"id"`
	Username       string         `json:"username"`
	VipLevel       NullVipType    `json:"vipLevel"`
	Totalbetamount int64          `json:"totalbetamount"`
	Newviplevel    string         `json:"newviplevel"`
	Rankprogress   pgtype.Numeric `json:"rankprogress"`
}

func (q *Queries) GetUserRanking(ctx context.Context, id pgtype.UUID) (GetUserRankingRow, error) {
	row := q.db.QueryRow(ctx, getUserRanking, id)
	var i GetUserRankingRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.VipLevel,
		&i.Totalbetamount,
		&i.Newviplevel,
		&i.Rankprogress,
	)
	return i, err
}

const getWeeklyTurnover = `-- name: GetWeeklyTurnover :one
SELECT
    COALESCE(SUM(b.bet), 0)::bigint AS weeklyBet
FROM
    "Bet" b
JOIN
    "User" u ON b."etgUsername" = u."etgUsername"
WHERE
    u.id = $1
    AND b."startTime" >= NOW() - INTERVAL '7 days'
    AND b."startTime" <= NOW()
`

func (q *Queries) GetWeeklyTurnover(ctx context.Context, id pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getWeeklyTurnover, id)
	var weeklybet int64
	err := row.Scan(&weeklybet)
	return weeklybet, err
}

const updateUserVipLevel = `-- name: UpdateUserVipLevel :exec
UPDATE
    "User"
SET
    "vipLevel" = (
        SELECT
            CASE
                WHEN COALESCE(SUM(b.bet), 0) >= 140000000 THEN 'KYAWTHUITE'
                WHEN COALESCE(SUM(b.bet), 0) >= 120000000 THEN 'JADE'
                WHEN COALESCE(SUM(b.bet), 0) >= 80000000 THEN 'DIAMOND_IV'
                WHEN COALESCE(SUM(b.bet), 0) >= 70000000 THEN 'DIAMOND_III'
                WHEN COALESCE(SUM(b.bet), 0) >= 60000000 THEN 'DIAMOND_II'
                WHEN COALESCE(SUM(b.bet), 0) >= 50000000 THEN 'DIAMOND_I'
                WHEN COALESCE(SUM(b.bet), 0) >= 40000000 THEN 'PLATINUM_IV'
                WHEN COALESCE(SUM(b.bet), 0) >= 35000000 THEN 'PLATINUM_III'
                WHEN COALESCE(SUM(b.bet), 0) >= 30000000 THEN 'PLATINUM_II'
                WHEN COALESCE(SUM(b.bet), 0) >= 25000000 THEN 'PLATINUM_I'
                WHEN COALESCE(SUM(b.bet), 0) >= 15000000 THEN 'GOLD'
                WHEN COALESCE(SUM(b.bet), 0) >= 7000000 THEN 'SILVER'
                ELSE 'BRONZE'
            END
        FROM
            "Bet" b
        WHERE
            b."etgUsername" = "User"."etgUsername"
    )
WHERE
    "User".id = $1
`

func (q *Queries) UpdateUserVipLevel(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, updateUserVipLevel, id)
	return err
}
