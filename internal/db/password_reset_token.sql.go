// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: password_reset_token.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPasswordResetToken = `-- name: CreatePasswordResetToken :exec
INSERT INTO "PasswordResetToken" ("tokenHash", "userId", "updatedAt") VALUES ($1, $2, now())
`

type CreatePasswordResetTokenParams struct {
	TokenHash string      `json:"tokenHash"`
	UserId    pgtype.UUID `json:"userId"`
}

func (q *Queries) CreatePasswordResetToken(ctx context.Context, arg CreatePasswordResetTokenParams) error {
	_, err := q.db.Exec(ctx, createPasswordResetToken, arg.TokenHash, arg.UserId)
	return err
}

const deletePasswordResetToken = `-- name: DeletePasswordResetToken :exec
DELETE FROM "PasswordResetToken" WHERE "tokenHash" = $1
`

func (q *Queries) DeletePasswordResetToken(ctx context.Context, tokenhash string) error {
	_, err := q.db.Exec(ctx, deletePasswordResetToken, tokenhash)
	return err
}

const getPasswordResetTokenByHash = `-- name: GetPasswordResetTokenByHash :one
SELECT prt."tokenHash", prt."userId", prt."createdAt", prt."updatedAt", u.username, u.email FROM "PasswordResetToken" prt JOIN "User" u ON prt."userId" = u.id WHERE "tokenHash" = $1
`

type GetPasswordResetTokenByHashRow struct {
	TokenHash string             `json:"tokenHash"`
	UserId    pgtype.UUID        `json:"userId"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
}

func (q *Queries) GetPasswordResetTokenByHash(ctx context.Context, tokenhash string) (GetPasswordResetTokenByHashRow, error) {
	row := q.db.QueryRow(ctx, getPasswordResetTokenByHash, tokenhash)
	var i GetPasswordResetTokenByHashRow
	err := row.Scan(
		&i.TokenHash,
		&i.UserId,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
	)
	return i, err
}

const getPasswordResetTokenByUserId = `-- name: GetPasswordResetTokenByUserId :one
SELECT prt."tokenHash", prt."userId", prt."createdAt", prt."updatedAt" FROM "PasswordResetToken" prt JOIN "User" u ON prt."userId" = u.id WHERE "userId" = $1
`

func (q *Queries) GetPasswordResetTokenByUserId(ctx context.Context, userid pgtype.UUID) (PasswordResetToken, error) {
	row := q.db.QueryRow(ctx, getPasswordResetTokenByUserId, userid)
	var i PasswordResetToken
	err := row.Scan(
		&i.TokenHash,
		&i.UserId,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
