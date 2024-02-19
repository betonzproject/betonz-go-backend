// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: bank.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createBank = `-- name: CreateBank :one
INSERT INTO
	"Bank" (id, "userId", name, "accountName", "accountNumber", "updatedAt")
VALUES
	(gen_random_uuid(), $1, $2, $3, $4, now())
RETURNING
	id, "userId", name, "accountName", "accountNumber", "createdAt", "updatedAt", disabled
`

type CreateBankParams struct {
	UserId        pgtype.UUID `json:"userId"`
	Name          BankName    `json:"name"`
	AccountName   string      `json:"accountName"`
	AccountNumber string      `json:"accountNumber"`
}

func (q *Queries) CreateBank(ctx context.Context, arg CreateBankParams) (Bank, error) {
	row := q.db.QueryRow(ctx, createBank,
		arg.UserId,
		arg.Name,
		arg.AccountName,
		arg.AccountNumber,
	)
	var i Bank
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Name,
		&i.AccountName,
		&i.AccountNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Disabled,
	)
	return i, err
}

const createSystemBank = `-- name: CreateSystemBank :one
INSERT INTO
	"Bank" (id, "userId", name, "accountName", "accountNumber", "updatedAt")
SELECT
	gen_random_uuid(),
	id,
	$1,
	$2,
	$3,
	now()
FROM
	"User"
WHERE
	role = 'SYSTEM'
LIMIT
	1
RETURNING
	id, "userId", name, "accountName", "accountNumber", "createdAt", "updatedAt", disabled
`

type CreateSystemBankParams struct {
	Name          BankName `json:"name"`
	AccountName   string   `json:"accountName"`
	AccountNumber string   `json:"accountNumber"`
}

func (q *Queries) CreateSystemBank(ctx context.Context, arg CreateSystemBankParams) (Bank, error) {
	row := q.db.QueryRow(ctx, createSystemBank, arg.Name, arg.AccountName, arg.AccountNumber)
	var i Bank
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Name,
		&i.AccountName,
		&i.AccountNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Disabled,
	)
	return i, err
}

const deleteBankById = `-- name: DeleteBankById :one
DELETE FROM "Bank" b USING "User" u WHERE b."userId" = u.id AND b.id = $1 RETURNING b.id, b."userId", b.name, b."accountName", b."accountNumber", b."createdAt", b."updatedAt", b.disabled
`

func (q *Queries) DeleteBankById(ctx context.Context, id pgtype.UUID) (Bank, error) {
	row := q.db.QueryRow(ctx, deleteBankById, id)
	var i Bank
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Name,
		&i.AccountName,
		&i.AccountNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Disabled,
	)
	return i, err
}

const deleteSystemBankById = `-- name: DeleteSystemBankById :one
DELETE FROM "Bank" b USING "User" u WHERE b."userId" = u.id AND b.id = $1 AND u.role = 'SYSTEM' RETURNING b.id, b."userId", b.name, b."accountName", b."accountNumber", b."createdAt", b."updatedAt", b.disabled
`

func (q *Queries) DeleteSystemBankById(ctx context.Context, id pgtype.UUID) (Bank, error) {
	row := q.db.QueryRow(ctx, deleteSystemBankById, id)
	var i Bank
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Name,
		&i.AccountName,
		&i.AccountNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Disabled,
	)
	return i, err
}

const getBankById = `-- name: GetBankById :one
SELECT id, "userId", name, "accountName", "accountNumber", "createdAt", "updatedAt", disabled FROM "Bank" WHERE id = $1
`

func (q *Queries) GetBankById(ctx context.Context, id pgtype.UUID) (Bank, error) {
	row := q.db.QueryRow(ctx, getBankById, id)
	var i Bank
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Name,
		&i.AccountName,
		&i.AccountNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Disabled,
	)
	return i, err
}

const getBanksByUserId = `-- name: GetBanksByUserId :many
SELECT id, "userId", name, "accountName", "accountNumber", "createdAt", "updatedAt", disabled FROM "Bank" WHERE "userId" = $1 ORDER BY "createdAt", "accountName"
`

func (q *Queries) GetBanksByUserId(ctx context.Context, userid pgtype.UUID) ([]Bank, error) {
	rows, err := q.db.Query(ctx, getBanksByUserId, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bank
	for rows.Next() {
		var i Bank
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.Name,
			&i.AccountName,
			&i.AccountNumber,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Disabled,
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

const getSystemBankById = `-- name: GetSystemBankById :one
SELECT b.id, b."userId", b.name, b."accountName", b."accountNumber", b."createdAt", b."updatedAt", b.disabled FROM "Bank" b JOIN "User" u ON b."userId" = u.id WHERE u.role = 'SYSTEM' AND b.id = $1
`

func (q *Queries) GetSystemBankById(ctx context.Context, id pgtype.UUID) (Bank, error) {
	row := q.db.QueryRow(ctx, getSystemBankById, id)
	var i Bank
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Name,
		&i.AccountName,
		&i.AccountNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Disabled,
	)
	return i, err
}

const getSystemBanks = `-- name: GetSystemBanks :many
SELECT b.id, b."userId", b.name, b."accountName", b."accountNumber", b."createdAt", b."updatedAt", b.disabled FROM "Bank" b JOIN "User" u ON b."userId" = u.id WHERE u.role = 'SYSTEM' ORDER BY b."createdAt", b."accountName"
`

func (q *Queries) GetSystemBanks(ctx context.Context) ([]Bank, error) {
	rows, err := q.db.Query(ctx, getSystemBanks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bank
	for rows.Next() {
		var i Bank
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.Name,
			&i.AccountName,
			&i.AccountNumber,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Disabled,
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

const getSystemBanksByBankName = `-- name: GetSystemBanksByBankName :many
SELECT b.id, b."userId", b.name, b."accountName", b."accountNumber", b."createdAt", b."updatedAt", b.disabled FROM "Bank" b JOIN "User" u ON b."userId" = u.id WHERE u.role = 'SYSTEM' AND b.name = $1 AND NOT disabled
`

func (q *Queries) GetSystemBanksByBankName(ctx context.Context, name BankName) ([]Bank, error) {
	rows, err := q.db.Query(ctx, getSystemBanksByBankName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bank
	for rows.Next() {
		var i Bank
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.Name,
			&i.AccountName,
			&i.AccountNumber,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Disabled,
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

const updateBank = `-- name: UpdateBank :exec
UPDATE "Bank"
SET
	name = $2,
	"accountName" = COALESCE($3, "accountName"),
	"accountNumber" = COALESCE($4, "accountNumber"),
	"updatedAt" = now()
WHERE
	id = $1
`

type UpdateBankParams struct {
	ID            pgtype.UUID `json:"id"`
	Name          BankName    `json:"name"`
	AccountName   pgtype.Text `json:"accountName"`
	AccountNumber pgtype.Text `json:"accountNumber"`
}

func (q *Queries) UpdateBank(ctx context.Context, arg UpdateBankParams) error {
	_, err := q.db.Exec(ctx, updateBank,
		arg.ID,
		arg.Name,
		arg.AccountName,
		arg.AccountNumber,
	)
	return err
}

const updateSystemBank = `-- name: UpdateSystemBank :one
UPDATE "Bank"
SET
	"accountName" = COALESCE($2, "accountName"),
	"accountNumber" = COALESCE($3, "accountNumber"),
	disabled = $4,
	"updatedAt" = now()
WHERE
	id = $1
RETURNING
	id, "userId", name, "accountName", "accountNumber", "createdAt", "updatedAt", disabled
`

type UpdateSystemBankParams struct {
	ID            pgtype.UUID `json:"id"`
	AccountName   pgtype.Text `json:"accountName"`
	AccountNumber pgtype.Text `json:"accountNumber"`
	Disabled      bool        `json:"disabled"`
}

func (q *Queries) UpdateSystemBank(ctx context.Context, arg UpdateSystemBankParams) (Bank, error) {
	row := q.db.QueryRow(ctx, updateSystemBank,
		arg.ID,
		arg.AccountName,
		arg.AccountNumber,
		arg.Disabled,
	)
	var i Bank
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Name,
		&i.AccountName,
		&i.AccountNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Disabled,
	)
	return i, err
}
