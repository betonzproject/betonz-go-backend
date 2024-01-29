// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getExtendedUserByUsername = `-- name: GetExtendedUserByUsername :one
SELECT id, username, email, "passwordHash", "displayName", "phoneNumber", "createdAt", "updatedAt", "etgUsername", role, "mainWallet", "lastUsedBankId", "profileImage", status, "lastLoginIp", "isEmailVerified", dob, "lastLoginAt", "pendingEmail" FROM "User" WHERE username = $1 AND ($2::"Role"[] IS NULL OR role = ANY($2))
`

type GetExtendedUserByUsernameParams struct {
	Username string `json:"username"`
	Roles    []Role `json:"roles"`
}

func (q *Queries) GetExtendedUserByUsername(ctx context.Context, arg GetExtendedUserByUsernameParams) (User, error) {
	row := q.db.QueryRow(ctx, getExtendedUserByUsername, arg.Username, arg.Roles)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.DisplayName,
		&i.PhoneNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.EtgUsername,
		&i.Role,
		&i.MainWallet,
		&i.LastUsedBankId,
		&i.ProfileImage,
		&i.Status,
		&i.LastLoginIp,
		&i.IsEmailVerified,
		&i.Dob,
		&i.LastLoginAt,
		&i.PendingEmail,
	)
	return i, err
}

const getPlayerInfoById = `-- name: GetPlayerInfoById :one
SELECT
	u.id,
	u.username,
	u.role,
	u.email,
	u.dob,
	u."displayName",
	u."phoneNumber",
	u."profileImage",
	u."mainWallet",
	u.status,
	u."createdAt",
	e."sourceIp" AS "lastLoginIp",
	e2."createdAt"::timestamptz AS "lastActiveAt"
FROM
	"User" u
LEFT JOIN (
	-- Get last login IP
	SELECT
		DISTINCT ON ("userId")
		"userId",
		"sourceIp"
	FROM
		"Event"
	WHERE
		result = 'SUCCESS'::"EventResult" AND type = 'LOGIN'::"EventType"
	ORDER BY
		"userId", "createdAt" DESC
) e ON 
	u.id = e."userId"
LEFT JOIN (
	-- Get last active time
	SELECT DISTINCT ON ("userId") "userId", "createdAt" FROM "Event" WHERE type = 'ACTIVE'::"EventType" ORDER BY "userId", "createdAt" DESC
) e2 ON
	u.id = e2."userId"
WHERE
	u.id = $1
`

type GetPlayerInfoByIdRow struct {
	ID           pgtype.UUID        `json:"id"`
	Username     string             `json:"username"`
	Role         Role               `json:"role"`
	Email        string             `json:"email"`
	Dob          pgtype.Date        `json:"dob"`
	DisplayName  pgtype.Text        `json:"displayName"`
	PhoneNumber  pgtype.Text        `json:"phoneNumber"`
	ProfileImage pgtype.Text        `json:"profileImage"`
	MainWallet   pgtype.Numeric     `json:"mainWallet"`
	Status       UserStatus         `json:"status"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	LastLoginIp  pgtype.Text        `json:"lastLoginIp"`
	LastActiveAt pgtype.Timestamptz `json:"lastActiveAt"`
}

func (q *Queries) GetPlayerInfoById(ctx context.Context, id pgtype.UUID) (GetPlayerInfoByIdRow, error) {
	row := q.db.QueryRow(ctx, getPlayerInfoById, id)
	var i GetPlayerInfoByIdRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.Email,
		&i.Dob,
		&i.DisplayName,
		&i.PhoneNumber,
		&i.ProfileImage,
		&i.MainWallet,
		&i.Status,
		&i.CreatedAt,
		&i.LastLoginIp,
		&i.LastActiveAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, username, role, email, "displayName", "phoneNumber", "mainWallet", dob, "profileImage", "lastUsedBankId", "isEmailVerified", status FROM "User" WHERE id = $1
`

type GetUserByIdRow struct {
	ID              pgtype.UUID    `json:"id"`
	Username        string         `json:"username"`
	Role            Role           `json:"role"`
	Email           string         `json:"email"`
	DisplayName     pgtype.Text    `json:"displayName"`
	PhoneNumber     pgtype.Text    `json:"phoneNumber"`
	MainWallet      pgtype.Numeric `json:"mainWallet"`
	Dob             pgtype.Date    `json:"dob"`
	ProfileImage    pgtype.Text    `json:"profileImage"`
	LastUsedBankId  pgtype.UUID    `json:"lastUsedBankId"`
	IsEmailVerified bool           `json:"isEmailVerified"`
	Status          UserStatus     `json:"status"`
}

func (q *Queries) GetUserById(ctx context.Context, id pgtype.UUID) (GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.Email,
		&i.DisplayName,
		&i.PhoneNumber,
		&i.MainWallet,
		&i.Dob,
		&i.ProfileImage,
		&i.LastUsedBankId,
		&i.IsEmailVerified,
		&i.Status,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT
	"rowNumber", id, username, role, email, dob, "displayName", "phoneNumber", "profileImage", "mainWallet", status, "createdAt", "lastLoginIp", "lastLogin", "lastDeposit", "lastWithdraw"
FROM (
	SELECT
		ROW_NUMBER() OVER (ORDER BY u."createdAt") "rowNumber",
		u.id,
		u.username,
		u.role,
		u.email,
		u.dob,
		u."displayName",
		u."phoneNumber",
		u."profileImage",
		u."mainWallet",
		u.status,
		u."createdAt",
		e."sourceIp" AS "lastLoginIp",
		e."createdAt"::timestamptz AS "lastLogin",
		tr1."lastDeposit"::timestamptz AS "lastDeposit",
		tr2."lastWithdraw"::timestamptz AS "lastWithdraw"
	FROM
		"User" u
	LEFT JOIN (
		-- Get last login IP and time
		SELECT
			DISTINCT ON ("userId")
			"userId",
			"sourceIp",
			"createdAt"
		FROM
			"Event"
		WHERE
			result = 'SUCCESS'::"EventResult" AND type = 'LOGIN'::"EventType"
		ORDER BY
			"userId", "createdAt" DESC
	) e ON 
		u.id = e."userId"
	LEFT JOIN (
		-- Get last deposit time
		SELECT
			"userId",
			max("updatedAt") "lastDeposit"
		FROM
			"TransactionRequest"
		WHERE
			type = 'DEPOSIT'::"TransactionType" AND status = 'APPROVED'::"TransactionStatus"
		GROUP BY
			"userId"
	) tr1 ON
		u.id = tr1."userId"
	LEFT JOIN (
		-- Get last withdraw time
		SELECT
			"userId",
			max("updatedAt") "lastWithdraw"
		FROM
			"TransactionRequest"
		WHERE
			type = 'WITHDRAW'::"TransactionType" AND status = 'APPROVED'::"TransactionStatus"
		GROUP BY
			"userId"
	) tr2 ON
		u.id = tr2."userId"
	WHERE
		u.role <> 'SYSTEM'::"Role"
	AND
		($1::"UserStatus"[] IS NULL OR u.status = ANY($1))
	AND (
		$2::text IS NULL
		OR u.username ILIKE '%' || $2 || '%' 
		OR u.email ILIKE '%' || $2 || '%'
		OR u."displayName" ILIKE '%' || $2 || '%'
		OR u."phoneNumber" ILIKE '%' || $2 || '%'
		OR u."lastLoginIp" ILIKE '%' || $2 || '%'
	)
	AND
		u."createdAt" >= $3
	AND
		u."createdAt" < $4
) q
ORDER BY
	"rowNumber" DESC, "createdAt" DESC
`

type GetUsersParams struct {
	Statuses []UserStatus       `json:"statuses"`
	Search   pgtype.Text        `json:"search"`
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

type GetUsersRow struct {
	RowNumber    int64              `json:"rowNumber"`
	ID           pgtype.UUID        `json:"id"`
	Username     string             `json:"username"`
	Role         Role               `json:"role"`
	Email        string             `json:"email"`
	Dob          pgtype.Date        `json:"dob"`
	DisplayName  pgtype.Text        `json:"displayName"`
	PhoneNumber  pgtype.Text        `json:"phoneNumber"`
	ProfileImage pgtype.Text        `json:"profileImage"`
	MainWallet   pgtype.Numeric     `json:"mainWallet"`
	Status       UserStatus         `json:"status"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	LastLoginIp  pgtype.Text        `json:"lastLoginIp"`
	LastLogin    pgtype.Timestamptz `json:"lastLogin"`
	LastDeposit  pgtype.Timestamptz `json:"lastDeposit"`
	LastWithdraw pgtype.Timestamptz `json:"lastWithdraw"`
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]GetUsersRow, error) {
	rows, err := q.db.Query(ctx, getUsers,
		arg.Statuses,
		arg.Search,
		arg.FromDate,
		arg.ToDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersRow
	for rows.Next() {
		var i GetUsersRow
		if err := rows.Scan(
			&i.RowNumber,
			&i.ID,
			&i.Username,
			&i.Role,
			&i.Email,
			&i.Dob,
			&i.DisplayName,
			&i.PhoneNumber,
			&i.ProfileImage,
			&i.MainWallet,
			&i.Status,
			&i.CreatedAt,
			&i.LastLoginIp,
			&i.LastLogin,
			&i.LastDeposit,
			&i.LastWithdraw,
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

const updateUser = `-- name: UpdateUser :exec
UPDATE "User" SET "displayName" = $2, email = $3, "phoneNumber" = $4, "updatedAt" = now() WHERE id = $1
`

type UpdateUserParams struct {
	ID          pgtype.UUID `json:"id"`
	DisplayName pgtype.Text `json:"displayName"`
	Email       string      `json:"email"`
	PhoneNumber pgtype.Text `json:"phoneNumber"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.DisplayName,
		arg.Email,
		arg.PhoneNumber,
	)
	return err
}

const updateUserLastUsedBank = `-- name: UpdateUserLastUsedBank :exec
UPDATE "User" SET "lastUsedBankId" = $2, "updatedAt" = now() WHERE id = $1
`

type UpdateUserLastUsedBankParams struct {
	ID             pgtype.UUID `json:"id"`
	LastUsedBankId pgtype.UUID `json:"lastUsedBankId"`
}

func (q *Queries) UpdateUserLastUsedBank(ctx context.Context, arg UpdateUserLastUsedBankParams) error {
	_, err := q.db.Exec(ctx, updateUserLastUsedBank, arg.ID, arg.LastUsedBankId)
	return err
}

const updateUserProfileImage = `-- name: UpdateUserProfileImage :exec
UPDATE "User" SET "profileImage" = $2, "updatedAt" = now() WHERE id = $1
`

type UpdateUserProfileImageParams struct {
	ID           pgtype.UUID `json:"id"`
	ProfileImage pgtype.Text `json:"profileImage"`
}

func (q *Queries) UpdateUserProfileImage(ctx context.Context, arg UpdateUserProfileImageParams) error {
	_, err := q.db.Exec(ctx, updateUserProfileImage, arg.ID, arg.ProfileImage)
	return err
}

const updateUserStatus = `-- name: UpdateUserStatus :exec
UPDATE "User" SET status = $2, "updatedAt" = now() WHERE id = $1
`

type UpdateUserStatusParams struct {
	ID     pgtype.UUID `json:"id"`
	Status UserStatus  `json:"status"`
}

func (q *Queries) UpdateUserStatus(ctx context.Context, arg UpdateUserStatusParams) error {
	_, err := q.db.Exec(ctx, updateUserStatus, arg.ID, arg.Status)
	return err
}
