// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: bet.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getBets = `-- name: GetBets :many
SELECT
	b.id, b."refId", b."etgUsername", b."providerUsername", b."productCode", b."productType", b."gameId", b.details, b.turnover, b.bet, b.payout, b.status, b."startTime", b."matchTime", b."endTime", b."settleTime", b."progShare", b."progWin", b.commission, b."winLoss",
	u.username,
	u.role,
	u."etgUsername"
FROM
	"Bet" b
	JOIN "User" u ON u."etgUsername" = b."etgUsername"
WHERE
	(
		$3::TEXT IS NULL
		OR b.id::TEXT ILIKE '%' || $3 || '%'
		OR u.username ILIKE '%' || $3 || '%'
		OR b."providerUsername" ILIKE '%' || $3 || '%'
	)
	AND (
		b."productCode" = $1
		OR $1 = 0
	)
	AND (
		b."productType" = $2
		OR $2 = 0
	)
	AND b."startTime" >= $4
	AND b."startTime" <= $5
ORDER BY
	b."startTime" DESC
`

type GetBetsParams struct {
	ProductCode int32              `json:"productCode"`
	ProductType int32              `json:"productType"`
	Search      pgtype.Text        `json:"search"`
	FromDate    pgtype.Timestamptz `json:"fromDate"`
	ToDate      pgtype.Timestamptz `json:"toDate"`
}

type GetBetsRow struct {
	ID               int32              `json:"id"`
	RefId            string             `json:"refId"`
	EtgUsername      string             `json:"etgUsername"`
	ProviderUsername string             `json:"providerUsername"`
	ProductCode      int32              `json:"productCode"`
	ProductType      int32              `json:"productType"`
	GameId           pgtype.Text        `json:"gameId"`
	Details          string             `json:"details"`
	Turnover         pgtype.Numeric     `json:"turnover"`
	Bet              pgtype.Numeric     `json:"bet"`
	Payout           pgtype.Numeric     `json:"payout"`
	Status           int32              `json:"status"`
	StartTime        pgtype.Timestamptz `json:"startTime"`
	MatchTime        pgtype.Timestamptz `json:"matchTime"`
	EndTime          pgtype.Timestamptz `json:"endTime"`
	SettleTime       pgtype.Timestamptz `json:"settleTime"`
	ProgShare        pgtype.Numeric     `json:"progShare"`
	ProgWin          pgtype.Numeric     `json:"progWin"`
	Commission       pgtype.Numeric     `json:"commission"`
	WinLoss          pgtype.Numeric     `json:"winLoss"`
	Username         string             `json:"username"`
	Role             Role               `json:"role"`
	EtgUsername_2    string             `json:"etgUsername_2"`
}

func (q *Queries) GetBets(ctx context.Context, arg GetBetsParams) ([]GetBetsRow, error) {
	rows, err := q.db.Query(ctx, getBets,
		arg.ProductCode,
		arg.ProductType,
		arg.Search,
		arg.FromDate,
		arg.ToDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetBetsRow{}
	for rows.Next() {
		var i GetBetsRow
		if err := rows.Scan(
			&i.ID,
			&i.RefId,
			&i.EtgUsername,
			&i.ProviderUsername,
			&i.ProductCode,
			&i.ProductType,
			&i.GameId,
			&i.Details,
			&i.Turnover,
			&i.Bet,
			&i.Payout,
			&i.Status,
			&i.StartTime,
			&i.MatchTime,
			&i.EndTime,
			&i.SettleTime,
			&i.ProgShare,
			&i.ProgWin,
			&i.Commission,
			&i.WinLoss,
			&i.Username,
			&i.Role,
			&i.EtgUsername_2,
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

const getTopPayout = `-- name: GetTopPayout :many
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
		$1::int IS NULL
		OR b1."productType" = $1
		OR $1 = 0
	)
	AND u.status = 'NORMAL'
`

type GetTopPayoutRow struct {
	ID           int32          `json:"id"`
	Username     string         `json:"username"`
	DisplayName  pgtype.Text    `json:"displayName"`
	ProfileImage pgtype.Text    `json:"profileImage"`
	Payout       pgtype.Numeric `json:"payout"`
}

func (q *Queries) GetTopPayout(ctx context.Context, producttype int32) ([]GetTopPayoutRow, error) {
	rows, err := q.db.Query(ctx, getTopPayout, producttype)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTopPayoutRow{}
	for rows.Next() {
		var i GetTopPayoutRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.DisplayName,
			&i.ProfileImage,
			&i.Payout,
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

const getTotalWinLoss = `-- name: GetTotalWinLoss :one
SELECT
	COALESCE(sum("winLoss"), 0)::bigint
FROM
	"Bet" b
	JOIN "User" u USING ("etgUsername")
WHERE
	u."role" = 'PLAYER'
	AND b."startTime" >= $1
	AND b."startTime" <= $2
`

type GetTotalWinLossParams struct {
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

func (q *Queries) GetTotalWinLoss(ctx context.Context, arg GetTotalWinLossParams) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalWinLoss, arg.FromDate, arg.ToDate)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getTurnoverByUserId = `-- name: GetTurnoverByUserId :many
SELECT
	b."productCode",
	sum(b.turnover) AS turnover
FROM
	"Bet" b
	JOIN "User" u USING ("etgUsername")
WHERE
	u.id = $1
	AND b."startTime" >= $2
	AND b."startTime" <= $3
GROUP BY
	b."productCode"
`

type GetTurnoverByUserIdParams struct {
	ID       pgtype.UUID        `json:"id"`
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

type GetTurnoverByUserIdRow struct {
	ProductCode int32 `json:"productCode"`
	Turnover    int64 `json:"turnover"`
}

func (q *Queries) GetTurnoverByUserId(ctx context.Context, arg GetTurnoverByUserIdParams) ([]GetTurnoverByUserIdRow, error) {
	rows, err := q.db.Query(ctx, getTurnoverByUserId, arg.ID, arg.FromDate, arg.ToDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTurnoverByUserIdRow{}
	for rows.Next() {
		var i GetTurnoverByUserIdRow
		if err := rows.Scan(&i.ProductCode, &i.Turnover); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertBet = `-- name: UpsertBet :exec
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
	"winLoss" = EXCLUDED."winLoss"
`

type UpsertBetParams struct {
	ID               int32              `json:"id"`
	RefId            string             `json:"refId"`
	EtgUsername      string             `json:"etgUsername"`
	ProviderUsername string             `json:"providerUsername"`
	ProductCode      int32              `json:"productCode"`
	ProductType      int32              `json:"productType"`
	GameId           pgtype.Text        `json:"gameId"`
	Details          string             `json:"details"`
	Turnover         pgtype.Numeric     `json:"turnover"`
	Bet              pgtype.Numeric     `json:"bet"`
	Payout           pgtype.Numeric     `json:"payout"`
	Status           int32              `json:"status"`
	StartTime        pgtype.Timestamptz `json:"startTime"`
	MatchTime        pgtype.Timestamptz `json:"matchTime"`
	EndTime          pgtype.Timestamptz `json:"endTime"`
	SettleTime       pgtype.Timestamptz `json:"settleTime"`
	ProgShare        pgtype.Numeric     `json:"progShare"`
	ProgWin          pgtype.Numeric     `json:"progWin"`
	Commission       pgtype.Numeric     `json:"commission"`
	WinLoss          pgtype.Numeric     `json:"winLoss"`
}

func (q *Queries) UpsertBet(ctx context.Context, arg UpsertBetParams) error {
	_, err := q.db.Exec(ctx, upsertBet,
		arg.ID,
		arg.RefId,
		arg.EtgUsername,
		arg.ProviderUsername,
		arg.ProductCode,
		arg.ProductType,
		arg.GameId,
		arg.Details,
		arg.Turnover,
		arg.Bet,
		arg.Payout,
		arg.Status,
		arg.StartTime,
		arg.MatchTime,
		arg.EndTime,
		arg.SettleTime,
		arg.ProgShare,
		arg.ProgWin,
		arg.Commission,
		arg.WinLoss,
	)
	return err
}
