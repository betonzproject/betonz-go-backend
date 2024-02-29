// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: identity_verification_request.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createIdentityVerificationRequest = `-- name: CreateIdentityVerificationRequest :exec
INSERT INTO
	"IdentityVerificationRequest" ("userId", "nricName", nric, dob, "nricFront", "nricBack", "holderFace", status, "updatedAt")
VALUES
	($1, $2, $3, $4, '', '', '', 'INCOMPLETE', now())
RETURNING
	id, "userId", "modifiedById", status, remarks, "nricFront", "nricBack", "holderFace", "nricName", nric, "createdAt", "updatedAt", dob
`

type CreateIdentityVerificationRequestParams struct {
	UserId   pgtype.UUID `json:"userId"`
	NricName string      `json:"nricName"`
	Nric     string      `json:"nric"`
	Dob      pgtype.Date `json:"dob"`
}

func (q *Queries) CreateIdentityVerificationRequest(ctx context.Context, arg CreateIdentityVerificationRequestParams) error {
	_, err := q.db.Exec(ctx, createIdentityVerificationRequest,
		arg.UserId,
		arg.NricName,
		arg.Nric,
		arg.Dob,
	)
	return err
}

const getIdentityVerificationRequestById = `-- name: GetIdentityVerificationRequestById :one
SELECT id, "userId", "modifiedById", status, remarks, "nricFront", "nricBack", "holderFace", "nricName", nric, "createdAt", "updatedAt", dob FROM "IdentityVerificationRequest" WHERE id = $1
`

func (q *Queries) GetIdentityVerificationRequestById(ctx context.Context, id int32) (IdentityVerificationRequest, error) {
	row := q.db.QueryRow(ctx, getIdentityVerificationRequestById, id)
	var i IdentityVerificationRequest
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.ModifiedById,
		&i.Status,
		&i.Remarks,
		&i.NricFront,
		&i.NricBack,
		&i.HolderFace,
		&i.NricName,
		&i.Nric,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Dob,
	)
	return i, err
}

const getIdentityVerificationRequests = `-- name: GetIdentityVerificationRequests :many
SELECT
	ivr.id, ivr."userId", ivr."modifiedById", ivr.status, ivr.remarks, ivr."nricFront", ivr."nricBack", ivr."holderFace", ivr."nricName", ivr.nric, ivr."createdAt", ivr."updatedAt", ivr.dob,
	u.username,
	u.role,
	u2.username AS "modifiedByUsername",
	u2.role AS "modifiedByRole"
FROM
	"IdentityVerificationRequest" ivr
	JOIN "User" u ON ivr."userId" = u.id
	LEFT JOIN "User" u2 ON ivr."modifiedById" = u2.id
WHERE
	(
		$1::"IdentityVerificationStatus"[] IS NULL
		OR ivr.status = ANY ($1)
	)
	AND (
		$2::TEXT IS NULL
		OR u.username ILIKE '%' || $2 || '%'
		OR u2.username ILIKE '%' || $2 || '%'
		OR ivr.id::TEXT ILIKE '%' || $2 || '%'
		OR ivr."nricName" ILIKE '%' || $2 || '%'
		OR ivr.nric ILIKE '%' || $2 || '%'
		OR ivr.remarks ILIKE '%' || $2 || '%'
	)
	AND ivr."createdAt" >= $3
	AND ivr."createdAt" <= $4
ORDER BY
	ivr.id DESC
`

type GetIdentityVerificationRequestsParams struct {
	Statuses []IdentityVerificationStatus `json:"statuses"`
	Search   pgtype.Text                  `json:"search"`
	FromDate pgtype.Timestamptz           `json:"fromDate"`
	ToDate   pgtype.Timestamptz           `json:"toDate"`
}

type GetIdentityVerificationRequestsRow struct {
	ID                 int32                      `json:"id"`
	UserId             pgtype.UUID                `json:"userId"`
	ModifiedById       pgtype.UUID                `json:"modifiedById"`
	Status             IdentityVerificationStatus `json:"status"`
	Remarks            pgtype.Text                `json:"remarks"`
	NricFront          string                     `json:"nricFront"`
	NricBack           string                     `json:"nricBack"`
	HolderFace         string                     `json:"holderFace"`
	NricName           string                     `json:"nricName"`
	Nric               string                     `json:"nric"`
	CreatedAt          pgtype.Timestamptz         `json:"createdAt"`
	UpdatedAt          pgtype.Timestamptz         `json:"updatedAt"`
	Dob                pgtype.Date                `json:"dob"`
	Username           string                     `json:"username"`
	Role               Role                       `json:"role"`
	ModifiedByUsername pgtype.Text                `json:"modifiedByUsername"`
	ModifiedByRole     NullRole                   `json:"modifiedByRole"`
}

func (q *Queries) GetIdentityVerificationRequests(ctx context.Context, arg GetIdentityVerificationRequestsParams) ([]GetIdentityVerificationRequestsRow, error) {
	rows, err := q.db.Query(ctx, getIdentityVerificationRequests,
		arg.Statuses,
		arg.Search,
		arg.FromDate,
		arg.ToDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetIdentityVerificationRequestsRow{}
	for rows.Next() {
		var i GetIdentityVerificationRequestsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.ModifiedById,
			&i.Status,
			&i.Remarks,
			&i.NricFront,
			&i.NricBack,
			&i.HolderFace,
			&i.NricName,
			&i.Nric,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Dob,
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

const getLatestIdentityVerificationRequestByUserId = `-- name: GetLatestIdentityVerificationRequestByUserId :one
SELECT id, "userId", "modifiedById", status, remarks, "nricFront", "nricBack", "holderFace", "nricName", nric, "createdAt", "updatedAt", dob FROM "IdentityVerificationRequest" WHERE "userId" = $1 AND status <> 'REJECTED' ORDER BY "createdAt" DESC LIMIT 1
`

func (q *Queries) GetLatestIdentityVerificationRequestByUserId(ctx context.Context, userid pgtype.UUID) (IdentityVerificationRequest, error) {
	row := q.db.QueryRow(ctx, getLatestIdentityVerificationRequestByUserId, userid)
	var i IdentityVerificationRequest
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.ModifiedById,
		&i.Status,
		&i.Remarks,
		&i.NricFront,
		&i.NricBack,
		&i.HolderFace,
		&i.NricName,
		&i.Nric,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Dob,
	)
	return i, err
}

const getPendingIdentityVerificationRequestCount = `-- name: GetPendingIdentityVerificationRequestCount :one
SELECT COUNT(*) FROM "IdentityVerificationRequest" WHERE status = 'PENDING'
`

func (q *Queries) GetPendingIdentityVerificationRequestCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getPendingIdentityVerificationRequestCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateIdentityVerificationRequestById = `-- name: UpdateIdentityVerificationRequestById :exec
UPDATE
	"IdentityVerificationRequest"
SET
	"modifiedById" = COALESCE($2, "modifiedById"),
	"nricName" = COALESCE($3, "nricName"),
	nric = COALESCE($4, nric),
	dob = COALESCE($5, dob),
	"nricFront" = COALESCE($6, "nricFront"),
	"nricBack" = COALESCE($7, "nricBack"),
	"holderFace" = COALESCE($8, "holderFace"),
	status = COALESCE($9, status),
	remarks = COALESCE($10, remarks),
	"updatedAt" = now()
WHERE
	id = $1
`

type UpdateIdentityVerificationRequestByIdParams struct {
	ID           int32                          `json:"id"`
	ModifiedById pgtype.UUID                    `json:"modifiedById"`
	NricName     pgtype.Text                    `json:"nricName"`
	Nric         pgtype.Text                    `json:"nric"`
	Dob          pgtype.Date                    `json:"dob"`
	NricFront    pgtype.Text                    `json:"nricFront"`
	NricBack     pgtype.Text                    `json:"nricBack"`
	HolderFace   pgtype.Text                    `json:"holderFace"`
	Status       NullIdentityVerificationStatus `json:"status"`
	Remarks      pgtype.Text                    `json:"remarks"`
}

func (q *Queries) UpdateIdentityVerificationRequestById(ctx context.Context, arg UpdateIdentityVerificationRequestByIdParams) error {
	_, err := q.db.Exec(ctx, updateIdentityVerificationRequestById,
		arg.ID,
		arg.ModifiedById,
		arg.NricName,
		arg.Nric,
		arg.Dob,
		arg.NricFront,
		arg.NricBack,
		arg.HolderFace,
		arg.Status,
		arg.Remarks,
	)
	return err
}
