// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
	"User" (id, username, email, "passwordHash", "etgUsername", "isEmailVerified", "updatedAt")
VALUES
	(gen_random_uuid (), $1, $2, $3, $4, true, now())
RETURNING
	id, username, email, "passwordHash", "displayName", "phoneNumber", "createdAt", "updatedAt", "etgUsername", role, "mainWallet", "lastUsedBankId", "profileImage", status, "isEmailVerified", dob, "pendingEmail"
`

type CreateUserParams struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	EtgUsername  string `json:"etgUsername"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.PasswordHash,
		arg.EtgUsername,
	)
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
		&i.IsEmailVerified,
		&i.Dob,
		&i.PendingEmail,
	)
	return i, err
}

const depositUserMainWallet = `-- name: DepositUserMainWallet :exec
UPDATE "User" SET "mainWallet" = "mainWallet" + $2, "updatedAt" = now() WHERE id = $1
`

type DepositUserMainWalletParams struct {
	ID     pgtype.UUID    `json:"id"`
	Amount pgtype.Numeric `json:"amount"`
}

func (q *Queries) DepositUserMainWallet(ctx context.Context, arg DepositUserMainWalletParams) error {
	_, err := q.db.Exec(ctx, depositUserMainWallet, arg.ID, arg.Amount)
	return err
}

const getExtendedUserById = `-- name: GetExtendedUserById :one
SELECT id, username, email, "passwordHash", "displayName", "phoneNumber", "createdAt", "updatedAt", "etgUsername", role, "mainWallet", "lastUsedBankId", "profileImage", status, "isEmailVerified", dob, "pendingEmail" FROM "User" WHERE id = $1
`

func (q *Queries) GetExtendedUserById(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getExtendedUserById, id)
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
		&i.IsEmailVerified,
		&i.Dob,
		&i.PendingEmail,
	)
	return i, err
}

const getExtendedUserByUsername = `-- name: GetExtendedUserByUsername :one
SELECT
	id, username, email, "passwordHash", "displayName", "phoneNumber", "createdAt", "updatedAt", "etgUsername", role, "mainWallet", "lastUsedBankId", "profileImage", status, "isEmailVerified", dob, "pendingEmail"
FROM
	"User"
WHERE
	username = $1
	AND (
		$2::"Role"[] IS NULL
		OR ROLE = ANY ($2)
	)
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
		&i.IsEmailVerified,
		&i.Dob,
		&i.PendingEmail,
	)
	return i, err
}

const getNewPlayerCount = `-- name: GetNewPlayerCount :one
SELECT
	COUNT(*)
FROM
	"User" u
WHERE
	u.role = 'PLAYER'
	AND u."createdAt" >= $1
	AND u."createdAt" <= $2
`

type GetNewPlayerCountParams struct {
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

func (q *Queries) GetNewPlayerCount(ctx context.Context, arg GetNewPlayerCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, getNewPlayerCount, arg.FromDate, arg.ToDate)
	var count int64
	err := row.Scan(&count)
	return count, err
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
	u."isEmailVerified",
	u."createdAt",
	e."sourceIp" AS "lastLoginIp",
	e2."updatedAt"::timestamptz AS "lastActiveAt"
FROM
	"User" u
	LEFT JOIN (
		-- Get last login IP
		SELECT DISTINCT
			ON ("userId") "userId",
			"sourceIp"
		FROM
			"Event"
		WHERE
			result = 'SUCCESS'
			AND type = 'LOGIN'
		ORDER BY
			"userId",
			"createdAt" DESC
	) e ON u.id = e."userId"
	LEFT JOIN (
		-- Get last active time
		SELECT DISTINCT
			ON ("userId") "userId",
			"updatedAt"
		FROM
			"Event"
		WHERE
			type = 'ACTIVE'
		ORDER BY
			"userId",
			"updatedAt" DESC
	) e2 ON u.id = e2."userId"
WHERE
	u.id = $1
`

type GetPlayerInfoByIdRow struct {
	ID              pgtype.UUID        `json:"id"`
	Username        string             `json:"username"`
	Role            Role               `json:"role"`
	Email           string             `json:"email"`
	Dob             pgtype.Date        `json:"dob"`
	DisplayName     pgtype.Text        `json:"displayName"`
	PhoneNumber     pgtype.Text        `json:"phoneNumber"`
	ProfileImage    pgtype.Text        `json:"profileImage"`
	MainWallet      pgtype.Numeric     `json:"mainWallet"`
	Status          UserStatus         `json:"status"`
	IsEmailVerified bool               `json:"isEmailVerified"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt"`
	LastLoginIp     pgtype.Text        `json:"lastLoginIp"`
	LastActiveAt    pgtype.Timestamptz `json:"lastActiveAt"`
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
		&i.IsEmailVerified,
		&i.CreatedAt,
		&i.LastLoginIp,
		&i.LastActiveAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, username, role, email, "pendingEmail", "displayName", "phoneNumber", "mainWallet", dob, "profileImage" FROM "User" WHERE id = $1
`

type GetUserByIdRow struct {
	ID           pgtype.UUID    `json:"id"`
	Username     string         `json:"username"`
	Role         Role           `json:"role"`
	Email        string         `json:"email"`
	PendingEmail pgtype.Text    `json:"pendingEmail"`
	DisplayName  pgtype.Text    `json:"displayName"`
	PhoneNumber  pgtype.Text    `json:"phoneNumber"`
	MainWallet   pgtype.Numeric `json:"mainWallet"`
	Dob          pgtype.Date    `json:"dob"`
	ProfileImage pgtype.Text    `json:"profileImage"`
}

func (q *Queries) GetUserById(ctx context.Context, id pgtype.UUID) (GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.Email,
		&i.PendingEmail,
		&i.DisplayName,
		&i.PhoneNumber,
		&i.MainWallet,
		&i.Dob,
		&i.ProfileImage,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
WITH q AS (
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
			SELECT DISTINCT
				ON ("userId") "userId",
				"sourceIp",
				"createdAt"
			FROM
				"Event"
			WHERE
				result = 'SUCCESS'
				AND type = 'LOGIN'
			ORDER BY
				"userId",
				"createdAt" DESC
		) e ON u.id = e."userId"
		LEFT JOIN (
			-- Get last deposit time
			SELECT
				"userId",
				max("updatedAt") "lastDeposit"
			FROM
				"TransactionRequest"
			WHERE
				type = 'DEPOSIT'
				AND status = 'APPROVED'
			GROUP BY
				"userId"
		) tr1 ON u.id = tr1."userId"
		LEFT JOIN (
			-- Get last withdraw time
			SELECT
				"userId",
				max("updatedAt") "lastWithdraw"
			FROM
				"TransactionRequest"
			WHERE
				type = 'WITHDRAW'
				AND status = 'APPROVED'
			GROUP BY
				"userId"
		) tr2 ON u.id = tr2."userId"
	WHERE
		u.role <> 'SYSTEM'
	ORDER BY
		u."createdAt"
)
SELECT
	"rowNumber", id, username, role, email, dob, "displayName", "phoneNumber", "profileImage", "mainWallet", status, "createdAt", "lastLoginIp", "lastLogin", "lastDeposit", "lastWithdraw"
FROM
	q
WHERE
	(
		$1::"UserStatus"[] IS NULL
		OR status = ANY ($1)
	)
	AND (
		$2::TEXT IS NULL
		OR username ILIKE '%' || $2 || '%'
		OR email ILIKE '%' || $2 || '%'
		OR "displayName" ILIKE '%' || $2 || '%'
		OR "phoneNumber" ILIKE '%' || $2 || '%'
		OR "lastLoginIp" ILIKE '%' || $2 || '%'
	)
	AND "createdAt" >= $3::timestamptz
	AND "createdAt" <= $4::timestamptz
ORDER BY
	"rowNumber" DESC
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
	items := []GetUsersRow{}
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

const markUserEmailAsVerified = `-- name: MarkUserEmailAsVerified :exec
UPDATE "User" SET "isEmailVerified" = true, email = $2, "updatedAt" = now() WHERE id = $1
`

type MarkUserEmailAsVerifiedParams struct {
	ID    pgtype.UUID `json:"id"`
	Email string      `json:"email"`
}

func (q *Queries) MarkUserEmailAsVerified(ctx context.Context, arg MarkUserEmailAsVerifiedParams) error {
	_, err := q.db.Exec(ctx, markUserEmailAsVerified, arg.ID, arg.Email)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE "User"
SET
	"displayName" = $2,
	"pendingEmail" = COALESCE($4, "pendingEmail"),
	"phoneNumber" = $3,
	"updatedAt" = now()
WHERE
	id = $1
`

type UpdateUserParams struct {
	ID           pgtype.UUID `json:"id"`
	DisplayName  pgtype.Text `json:"displayName"`
	PhoneNumber  pgtype.Text `json:"phoneNumber"`
	PendingEmail pgtype.Text `json:"pendingEmail"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.DisplayName,
		arg.PhoneNumber,
		arg.PendingEmail,
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

const updateUserPasswordHash = `-- name: UpdateUserPasswordHash :exec
UPDATE "User" SET "passwordHash" = $2, "updatedAt" = now() WHERE id = $1
`

type UpdateUserPasswordHashParams struct {
	ID           pgtype.UUID `json:"id"`
	PasswordHash string      `json:"passwordHash"`
}

func (q *Queries) UpdateUserPasswordHash(ctx context.Context, arg UpdateUserPasswordHashParams) error {
	_, err := q.db.Exec(ctx, updateUserPasswordHash, arg.ID, arg.PasswordHash)
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

const updateUsername = `-- name: UpdateUsername :exec
UPDATE "User" SET username = $2, "updatedAt" = now() WHERE id = $1
`

type UpdateUsernameParams struct {
	ID       pgtype.UUID `json:"id"`
	Username string      `json:"username"`
}

func (q *Queries) UpdateUsername(ctx context.Context, arg UpdateUsernameParams) error {
	_, err := q.db.Exec(ctx, updateUsername, arg.ID, arg.Username)
	return err
}

const withdrawUserMainWallet = `-- name: WithdrawUserMainWallet :exec
UPDATE "User" SET "mainWallet" = "mainWallet" - $2, "updatedAt" = now() WHERE id = $1
`

type WithdrawUserMainWalletParams struct {
	ID     pgtype.UUID    `json:"id"`
	Amount pgtype.Numeric `json:"amount"`
}

func (q *Queries) WithdrawUserMainWallet(ctx context.Context, arg WithdrawUserMainWalletParams) error {
	_, err := q.db.Exec(ctx, withdrawUserMainWallet, arg.ID, arg.Amount)
	return err
}
