// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: transaction_log.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getTransactionLogs = `-- name: GetTransactionLogs :many
SELECT
	tr."initiatorId", tr."beneficiaryId", tr.product, tr."balanceBefore", tr."balanceAfter", tr.amount, tr.type, tr.remarks, tr."createdAt", tr."updatedAt", tr.id, tr."receiptPath", tr.bonus,
	u.username AS "initiatorUsername",
	u.role AS "initiatorRole",
	u2.username AS "beneficiaryUsername",
	u2.role AS "beneficiaryRole"
FROM
	"Transaction" tr
	JOIN "User" u ON u.id = tr."initiatorId"
	JOIN "User" u2 ON u2.id = tr."beneficiaryId"
WHERE
	(
		$1::TEXT IS NULL
		OR u.username ILIKE '%' || $1 || '%'
		OR u2.username ILIKE '%' || $1 || '%'
		OR tr.remarks ILIKE '%' || $1 || '%'
	)
	AND tr."createdAt" >= $2
	AND tr."createdAt" <= $3
ORDER BY
	tr.id DESC
`

type GetTransactionLogsParams struct {
	Search   pgtype.Text        `json:"search"`
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

type GetTransactionLogsRow struct {
	InitiatorId         pgtype.UUID        `json:"initiatorId"`
	BeneficiaryId       pgtype.UUID        `json:"beneficiaryId"`
	Product             string             `json:"product"`
	BalanceBefore       pgtype.Numeric     `json:"balanceBefore"`
	BalanceAfter        pgtype.Numeric     `json:"balanceAfter"`
	Amount              pgtype.Numeric     `json:"amount"`
	Type                TransactionType    `json:"type"`
	Remarks             pgtype.Text        `json:"remarks"`
	CreatedAt           pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt           pgtype.Timestamptz `json:"updatedAt"`
	ID                  int32              `json:"id"`
	ReceiptPath         pgtype.Text        `json:"receiptPath"`
	Bonus               pgtype.Numeric     `json:"bonus"`
	InitiatorUsername   string             `json:"initiatorUsername"`
	InitiatorRole       Role               `json:"initiatorRole"`
	BeneficiaryUsername string             `json:"beneficiaryUsername"`
	BeneficiaryRole     Role               `json:"beneficiaryRole"`
}

func (q *Queries) GetTransactionLogs(ctx context.Context, arg GetTransactionLogsParams) ([]GetTransactionLogsRow, error) {
	rows, err := q.db.Query(ctx, getTransactionLogs, arg.Search, arg.FromDate, arg.ToDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransactionLogsRow
	for rows.Next() {
		var i GetTransactionLogsRow
		if err := rows.Scan(
			&i.InitiatorId,
			&i.BeneficiaryId,
			&i.Product,
			&i.BalanceBefore,
			&i.BalanceAfter,
			&i.Amount,
			&i.Type,
			&i.Remarks,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID,
			&i.ReceiptPath,
			&i.Bonus,
			&i.InitiatorUsername,
			&i.InitiatorRole,
			&i.BeneficiaryUsername,
			&i.BeneficiaryRole,
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
