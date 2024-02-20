// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: transaction_request.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTransactionRequest = `-- name: CreateTransactionRequest :exec
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
		"updatedAt"
	)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, now())
`

type CreateTransactionRequestParams struct {
	UserId                       pgtype.UUID       `json:"userId"`
	BankName                     BankName          `json:"bankName"`
	BankAccountName              string            `json:"bankAccountName"`
	BankAccountNumber            string            `json:"bankAccountNumber"`
	BeneficiaryBankAccountName   pgtype.Text       `json:"beneficiaryBankAccountName"`
	BeneficiaryBankAccountNumber pgtype.Text       `json:"beneficiaryBankAccountNumber"`
	Amount                       pgtype.Numeric    `json:"amount"`
	Bonus                        pgtype.Numeric    `json:"bonus"`
	Type                         TransactionType   `json:"type"`
	DepositToWallet              pgtype.Int4       `json:"depositToWallet"`
	Promotion                    NullPromotionType `json:"promotion"`
	ReceiptPath                  pgtype.Text       `json:"receiptPath"`
	Status                       TransactionStatus `json:"status"`
}

func (q *Queries) CreateTransactionRequest(ctx context.Context, arg CreateTransactionRequestParams) error {
	_, err := q.db.Exec(ctx, createTransactionRequest,
		arg.UserId,
		arg.BankName,
		arg.BankAccountName,
		arg.BankAccountNumber,
		arg.BeneficiaryBankAccountName,
		arg.BeneficiaryBankAccountNumber,
		arg.Amount,
		arg.Bonus,
		arg.Type,
		arg.DepositToWallet,
		arg.Promotion,
		arg.ReceiptPath,
		arg.Status,
	)
	return err
}

const getNewPlayerWithTransactionsCount = `-- name: GetNewPlayerWithTransactionsCount :one
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
			AND tr."updatedAt" >= $1
			AND tr."updatedAt" <= $2
			AND u."createdAt" >= $1
			AND u."createdAd" <= $2
	) q
`

type GetNewPlayerWithTransactionsCountParams struct {
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

func (q *Queries) GetNewPlayerWithTransactionsCount(ctx context.Context, arg GetNewPlayerWithTransactionsCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, getNewPlayerWithTransactionsCount, arg.FromDate, arg.ToDate)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getPlayerWithTransactionsCount = `-- name: GetPlayerWithTransactionsCount :one
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
			AND "updatedAt" >= $1
			AND "updatedAt" <= $2
	) q
`

type GetPlayerWithTransactionsCountParams struct {
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

func (q *Queries) GetPlayerWithTransactionsCount(ctx context.Context, arg GetPlayerWithTransactionsCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, getPlayerWithTransactionsCount, arg.FromDate, arg.ToDate)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalTransactionAmountAndCount = `-- name: GetTotalTransactionAmountAndCount :one
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
	AND tr."updatedAt" >= $2
	AND tr."updatedAt" <= $3
`

type GetTotalTransactionAmountAndCountParams struct {
	Type     TransactionType    `json:"type"`
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

type GetTotalTransactionAmountAndCountRow struct {
	Total      int64 `json:"total"`
	Count      int64 `json:"count"`
	BonusTotal int64 `json:"bonusTotal"`
}

func (q *Queries) GetTotalTransactionAmountAndCount(ctx context.Context, arg GetTotalTransactionAmountAndCountParams) (GetTotalTransactionAmountAndCountRow, error) {
	row := q.db.QueryRow(ctx, getTotalTransactionAmountAndCount, arg.Type, arg.FromDate, arg.ToDate)
	var i GetTotalTransactionAmountAndCountRow
	err := row.Scan(&i.Total, &i.Count, &i.BonusTotal)
	return i, err
}

const getTransactionRequestById = `-- name: GetTransactionRequestById :one
SELECT id, "userId", "modifiedById", "bankName", "bankAccountName", "bankAccountNumber", "beneficiaryBankAccountName", "beneficiaryBankAccountNumber", amount, type, "receiptPath", status, remarks, "createdAt", "updatedAt", bonus, "withdrawBankFees", "depositToWallet", promotion FROM "TransactionRequest" WHERE id = $1
`

func (q *Queries) GetTransactionRequestById(ctx context.Context, id int32) (TransactionRequest, error) {
	row := q.db.QueryRow(ctx, getTransactionRequestById, id)
	var i TransactionRequest
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.ModifiedById,
		&i.BankName,
		&i.BankAccountName,
		&i.BankAccountNumber,
		&i.BeneficiaryBankAccountName,
		&i.BeneficiaryBankAccountNumber,
		&i.Amount,
		&i.Type,
		&i.ReceiptPath,
		&i.Status,
		&i.Remarks,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Bonus,
		&i.WithdrawBankFees,
		&i.DepositToWallet,
		&i.Promotion,
	)
	return i, err
}

const getTransactionRequests = `-- name: GetTransactionRequests :many
SELECT
	tr.id, tr."userId", tr."modifiedById", tr."bankName", tr."bankAccountName", tr."bankAccountNumber", tr."beneficiaryBankAccountName", tr."beneficiaryBankAccountNumber", tr.amount, tr.type, tr."receiptPath", tr.status, tr.remarks, tr."createdAt", tr."updatedAt", tr.bonus, tr."withdrawBankFees", tr."depositToWallet", tr.promotion,
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
		$1::"TransactionType"[] IS NULL
		OR tr."type" = ANY ($1)
	)
	AND (
		$2::"TransactionStatus"[] IS NULL
		OR tr.status = ANY ($2)
	)
	AND (
		$3::TEXT IS NULL
		OR u.username ILIKE '%' || $3 || '%'
		OR u2.username ILIKE '%' || $3 || '%'
		OR tr."bankAccountName" ILIKE '%' || $3 || '%'
		OR tr."beneficiaryBankAccountName" ILIKE '%' || $3 || '%'
		OR tr.remarks ILIKE '%' || $3 || '%'
	)
	AND tr."createdAt" >= $4
	AND tr."createdAt" <= $5
ORDER BY
	tr.id DESC
`

type GetTransactionRequestsParams struct {
	Types    []TransactionType   `json:"types"`
	Statuses []TransactionStatus `json:"statuses"`
	Search   pgtype.Text         `json:"search"`
	FromDate pgtype.Timestamptz  `json:"fromDate"`
	ToDate   pgtype.Timestamptz  `json:"toDate"`
}

type GetTransactionRequestsRow struct {
	ID                           int32              `json:"id"`
	UserId                       pgtype.UUID        `json:"userId"`
	ModifiedById                 pgtype.UUID        `json:"modifiedById"`
	BankName                     BankName           `json:"bankName"`
	BankAccountName              string             `json:"bankAccountName"`
	BankAccountNumber            string             `json:"bankAccountNumber"`
	BeneficiaryBankAccountName   pgtype.Text        `json:"beneficiaryBankAccountName"`
	BeneficiaryBankAccountNumber pgtype.Text        `json:"beneficiaryBankAccountNumber"`
	Amount                       pgtype.Numeric     `json:"amount"`
	Type                         TransactionType    `json:"type"`
	ReceiptPath                  pgtype.Text        `json:"receiptPath"`
	Status                       TransactionStatus  `json:"status"`
	Remarks                      pgtype.Text        `json:"remarks"`
	CreatedAt                    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt                    pgtype.Timestamptz `json:"updatedAt"`
	Bonus                        pgtype.Numeric     `json:"bonus"`
	WithdrawBankFees             pgtype.Numeric     `json:"withdrawBankFees"`
	DepositToWallet              pgtype.Int4        `json:"depositToWallet"`
	Promotion                    NullPromotionType  `json:"promotion"`
	Username                     string             `json:"username"`
	Role                         Role               `json:"role"`
	ModifiedByUsername           pgtype.Text        `json:"modifiedByUsername"`
	ModifiedByRole               NullRole           `json:"modifiedByRole"`
}

func (q *Queries) GetTransactionRequests(ctx context.Context, arg GetTransactionRequestsParams) ([]GetTransactionRequestsRow, error) {
	rows, err := q.db.Query(ctx, getTransactionRequests,
		arg.Types,
		arg.Statuses,
		arg.Search,
		arg.FromDate,
		arg.ToDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTransactionRequestsRow{}
	for rows.Next() {
		var i GetTransactionRequestsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.ModifiedById,
			&i.BankName,
			&i.BankAccountName,
			&i.BankAccountNumber,
			&i.BeneficiaryBankAccountName,
			&i.BeneficiaryBankAccountNumber,
			&i.Amount,
			&i.Type,
			&i.ReceiptPath,
			&i.Status,
			&i.Remarks,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Bonus,
			&i.WithdrawBankFees,
			&i.DepositToWallet,
			&i.Promotion,
			&i.Username,
			&i.Role,
			&i.ModifiedByUsername,
			&i.ModifiedByRole,
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

const getTransactionRequestsByUserId = `-- name: GetTransactionRequestsByUserId :many
SELECT
	tr.id, tr."userId", tr."modifiedById", tr."bankName", tr."bankAccountName", tr."bankAccountNumber", tr."beneficiaryBankAccountName", tr."beneficiaryBankAccountNumber", tr.amount, tr.type, tr."receiptPath", tr.status, tr.remarks, tr."createdAt", tr."updatedAt", tr.bonus, tr."withdrawBankFees", tr."depositToWallet", tr.promotion,
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
		$2::"TransactionType"[] IS NULL
		OR tr."type" = ANY ($2)
	)
	AND (
		$3::"TransactionStatus"[] IS NULL
		OR tr.status = ANY ($3)
	)
	AND tr."createdAt" >= $4
	AND tr."createdAt" <= $5
	AND tr."userId" = $1
ORDER BY
	tr.id DESC
`

type GetTransactionRequestsByUserIdParams struct {
	UserId   pgtype.UUID         `json:"userId"`
	Types    []TransactionType   `json:"types"`
	Statuses []TransactionStatus `json:"statuses"`
	FromDate pgtype.Timestamptz  `json:"fromDate"`
	ToDate   pgtype.Timestamptz  `json:"toDate"`
}

type GetTransactionRequestsByUserIdRow struct {
	ID                           int32              `json:"id"`
	UserId                       pgtype.UUID        `json:"userId"`
	ModifiedById                 pgtype.UUID        `json:"modifiedById"`
	BankName                     BankName           `json:"bankName"`
	BankAccountName              string             `json:"bankAccountName"`
	BankAccountNumber            string             `json:"bankAccountNumber"`
	BeneficiaryBankAccountName   pgtype.Text        `json:"beneficiaryBankAccountName"`
	BeneficiaryBankAccountNumber pgtype.Text        `json:"beneficiaryBankAccountNumber"`
	Amount                       pgtype.Numeric     `json:"amount"`
	Type                         TransactionType    `json:"type"`
	ReceiptPath                  pgtype.Text        `json:"receiptPath"`
	Status                       TransactionStatus  `json:"status"`
	Remarks                      pgtype.Text        `json:"remarks"`
	CreatedAt                    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt                    pgtype.Timestamptz `json:"updatedAt"`
	Bonus                        pgtype.Numeric     `json:"bonus"`
	WithdrawBankFees             pgtype.Numeric     `json:"withdrawBankFees"`
	DepositToWallet              pgtype.Int4        `json:"depositToWallet"`
	Promotion                    NullPromotionType  `json:"promotion"`
	Username                     string             `json:"username"`
	Role                         Role               `json:"role"`
	ModifiedByUsername           pgtype.Text        `json:"modifiedByUsername"`
	ModifiedByRole               NullRole           `json:"modifiedByRole"`
}

func (q *Queries) GetTransactionRequestsByUserId(ctx context.Context, arg GetTransactionRequestsByUserIdParams) ([]GetTransactionRequestsByUserIdRow, error) {
	rows, err := q.db.Query(ctx, getTransactionRequestsByUserId,
		arg.UserId,
		arg.Types,
		arg.Statuses,
		arg.FromDate,
		arg.ToDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTransactionRequestsByUserIdRow{}
	for rows.Next() {
		var i GetTransactionRequestsByUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.ModifiedById,
			&i.BankName,
			&i.BankAccountName,
			&i.BankAccountNumber,
			&i.BeneficiaryBankAccountName,
			&i.BeneficiaryBankAccountNumber,
			&i.Amount,
			&i.Type,
			&i.ReceiptPath,
			&i.Status,
			&i.Remarks,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Bonus,
			&i.WithdrawBankFees,
			&i.DepositToWallet,
			&i.Promotion,
			&i.Username,
			&i.Role,
			&i.ModifiedByUsername,
			&i.ModifiedByRole,
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

const hasApprovedDepositRequestsWithin30DaysByUserId = `-- name: HasApprovedDepositRequestsWithin30DaysByUserId :one
SELECT
	EXISTS (
		SELECT
			id, "userId", "modifiedById", "bankName", "bankAccountName", "bankAccountNumber", "beneficiaryBankAccountName", "beneficiaryBankAccountNumber", amount, type, "receiptPath", status, remarks, "createdAt", "updatedAt", bonus, "withdrawBankFees", "depositToWallet", promotion
		FROM
			"TransactionRequest"
		WHERE
			"userId" = $1
			AND type = 'DEPOSIT'
			AND status = 'APPROVED'
			AND "updatedAt" >= now() - INTERVAL '30 days'
	)
`

func (q *Queries) HasApprovedDepositRequestsWithin30DaysByUserId(ctx context.Context, userid pgtype.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, hasApprovedDepositRequestsWithin30DaysByUserId, userid)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const hasPendingTransactionRequestsWithPromotion = `-- name: HasPendingTransactionRequestsWithPromotion :one
SELECT
	EXISTS (
		SELECT
			id, "userId", "modifiedById", "bankName", "bankAccountName", "bankAccountNumber", "beneficiaryBankAccountName", "beneficiaryBankAccountNumber", amount, type, "receiptPath", status, remarks, "createdAt", "updatedAt", bonus, "withdrawBankFees", "depositToWallet", promotion
		FROM
			"TransactionRequest"
		WHERE
			"userId" = $1
			AND type = 'DEPOSIT'
			AND status = 'PENDING'
			AND promotion = $2
	)
`

type HasPendingTransactionRequestsWithPromotionParams struct {
	UserId    pgtype.UUID       `json:"userId"`
	Promotion NullPromotionType `json:"promotion"`
}

func (q *Queries) HasPendingTransactionRequestsWithPromotion(ctx context.Context, arg HasPendingTransactionRequestsWithPromotionParams) (bool, error) {
	row := q.db.QueryRow(ctx, hasPendingTransactionRequestsWithPromotion, arg.UserId, arg.Promotion)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const hasRecentDepositRequestsByUserId = `-- name: HasRecentDepositRequestsByUserId :one
SELECT
	EXISTS (
		SELECT
			id, "userId", "modifiedById", "bankName", "bankAccountName", "bankAccountNumber", "beneficiaryBankAccountName", "beneficiaryBankAccountNumber", amount, type, "receiptPath", status, remarks, "createdAt", "updatedAt", bonus, "withdrawBankFees", "depositToWallet", promotion
		FROM
			"TransactionRequest"
		WHERE
			"userId" = $1
			AND type = 'DEPOSIT'
			AND status = 'PENDING'
			AND "createdAt" >= now() - INTERVAL '1 minute'
	)
`

func (q *Queries) HasRecentDepositRequestsByUserId(ctx context.Context, userid pgtype.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, hasRecentDepositRequestsByUserId, userid)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const hasRecentWithdrawRequestsByUserId = `-- name: HasRecentWithdrawRequestsByUserId :one
SELECT
	EXISTS (
		SELECT
			id, "userId", "modifiedById", "bankName", "bankAccountName", "bankAccountNumber", "beneficiaryBankAccountName", "beneficiaryBankAccountNumber", amount, type, "receiptPath", status, remarks, "createdAt", "updatedAt", bonus, "withdrawBankFees", "depositToWallet", promotion
		FROM
			"TransactionRequest"
		WHERE
			"userId" = $1
			AND type = 'WITHDRAW'
			AND status = 'PENDING'
			AND "createdAt" >= now() - INTERVAL '5 minutes'
	)
`

func (q *Queries) HasRecentWithdrawRequestsByUserId(ctx context.Context, userid pgtype.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, hasRecentWithdrawRequestsByUserId, userid)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const updateTransactionRequest = `-- name: UpdateTransactionRequest :exec
UPDATE "TransactionRequest"
SET
	"modifiedById" = $2,
	"receiptPath" = COALESCE($3, "receiptPath"),
	status = $4,
	"withdrawBankFees" = COALESCE($5, 0),
	remarks = $6,
	"updatedAt" = now()
WHERE id = $1
`

type UpdateTransactionRequestParams struct {
	ID               int32             `json:"id"`
	ModifiedById     pgtype.UUID       `json:"modifiedById"`
	ReceiptPath      pgtype.Text       `json:"receiptPath"`
	Status           TransactionStatus `json:"status"`
	WithdrawBankFees pgtype.Numeric    `json:"withdrawBankFees"`
	Remarks          pgtype.Text       `json:"remarks"`
}

func (q *Queries) UpdateTransactionRequest(ctx context.Context, arg UpdateTransactionRequestParams) error {
	_, err := q.db.Exec(ctx, updateTransactionRequest,
		arg.ID,
		arg.ModifiedById,
		arg.ReceiptPath,
		arg.Status,
		arg.WithdrawBankFees,
		arg.Remarks,
	)
	return err
}
