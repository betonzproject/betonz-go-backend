-- name: GetUserById :one
SELECT
	id,
	username,
	role,
	email,
	"pendingEmail",
	"displayName",
	"phoneNumber",
	"mainWallet",
	dob,
	"referralCode",
	"level",
	"exp",
	"betonPoint",
	"profileImage",
	"vipLevel"
FROM
	"User"
WHERE
	id = $1;

-- name: GetExtendedUserById :one
SELECT
	*
FROM
	"User"
WHERE
	id = $1;

-- name: GetExtendedUserByUsername :one
SELECT
	*
FROM
	"User"
WHERE
	username = $1
	AND (
		@roles::"Role"[] IS NULL
		OR role = ANY (@roles)
	);

-- name: GetUsers :many
WITH
	q AS (
		SELECT
			ROW_NUMBER() OVER (
				ORDER BY
					u."createdAt"
			) "rowNumber",
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
			u."referralCode",
			u."createdAt",
			(
				SELECT
					COUNT(*)
				FROM
					"User" u2
				WHERE
					u."referralCode" = u2."invitedBy"
			) AS "invitedUserCount",
			e."sourceIp" AS "lastLoginIp",
			e."createdAt"::timestamptz AS "lastLogin",
			tr1."lastDeposit"::timestamptz AS "lastDeposit",
			tr2."lastWithdraw"::timestamptz AS "lastWithdraw",
			u."vipLevel"
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
	*
FROM
	q
WHERE
	(
		@statuses::"UserStatus"[] IS NULL
		OR status = ANY (@statuses)
	)
	AND (
		sqlc.narg('search')::TEXT IS NULL
		OR "rowNumber"::TEXT ILIKE '%' || @search || '%'
		OR username ILIKE '%' || @search || '%'
		OR email ILIKE '%' || @search || '%'
		OR "displayName" ILIKE '%' || @search || '%'
		OR "phoneNumber" ILIKE '%' || @search || '%'
		OR "lastLoginIp" ILIKE '%' || @search || '%'
	)
	AND "createdAt" >= sqlc.arg('fromDate')::timestamptz
	AND "createdAt" <= sqlc.arg('toDate')::timestamptz
ORDER BY
	"rowNumber" DESC;

-- name: GetPlayerByReferralCode :one
SELECT
	*
FROM
	"User"
WHERE
	"referralCode" = $1;

-- name: GetPlayerInfoById :one
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
	u."referralCode",
	u."isEmailVerified",
	u."createdAt",
	e."sourceIp" AS "lastLoginIp",
	e2."updatedAt"::timestamptz AS "lastActiveAt",
	u."vipLevel"
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
	u.id = $1;

-- name: GetNewPlayerCount :one
SELECT
	COUNT(*)
FROM
	"User" u
WHERE
	u.role = 'PLAYER'
	AND u."createdAt" >= sqlc.arg('fromDate')
	AND u."createdAt" <= sqlc.arg('toDate');

-- name: CreateUser :one
INSERT INTO
	"User" (
		id,
		username,
		email,
		"passwordHash",
		"etgUsername",
		"isEmailVerified",
		"referralCode",
		"invitedBy",
		"updatedAt"
	)
VALUES
	(gen_random_uuid (), $1, $2, $3, $4, TRUE, $5, $6, now())
RETURNING
	*;

-- name: CreateAdmin :one
INSERT INTO
	"User" (id, username, email, "etgUsername", "passwordHash", "isEmailVerified", role, "updatedAt")
VALUES
	(gen_random_uuid (), $1, $2, $3, $4, TRUE, $5, now())
RETURNING
	*;

-- name: GetAdmins :many
WITH
	q AS (
		SELECT
			ROW_NUMBER() OVER (
				ORDER BY
					u."createdAt"
			) "rowNumber",
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
			u."referralCode",
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
	*
FROM
	q
WHERE
	(
		@role::"Role"[] IS NULL
		OR role = ANY (@role)
	)
ORDER BY
	"rowNumber" DESC;

-- name: DeleteAdmin :exec
DELETE FROM "User"
WHERE
	id = $1
	AND (
		role = 'ADMIN'
		OR role = 'SUPERADMIN'
	);

-- name: UpdateUser :exec
UPDATE "User"
SET
	"displayName" = $2,
	"pendingEmail" = COALESCE(sqlc.arg('pendingEmail'), "pendingEmail"),
	"phoneNumber" = $3,
	"isEmailVerified" = $4,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: UpdateUsername :exec
UPDATE "User"
SET
	username = $2,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: UpdateUserPasswordHash :exec
UPDATE "User"
SET
	"passwordHash" = $2,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: UpdateUserDob :exec
UPDATE "User"
SET
	dob = $2,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: UpdateUserProfileImage :exec
UPDATE "User"
SET
	"profileImage" = $2,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: DepositUserMainWallet :exec
UPDATE "User"
SET
	"mainWallet" = "mainWallet" + @amount,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: WithdrawUserMainWallet :exec
UPDATE "User"
SET
	"mainWallet" = "mainWallet" - @amount,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: UpdateUserLastUsedBank :exec
UPDATE "User"
SET
	"lastUsedBankId" = $2,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: UpdateUserStatus :exec
UPDATE "User"
SET
	status = $2,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: MarkUserEmailAsVerified :exec
UPDATE "User"
SET
	"isEmailVerified" = TRUE,
	email = $2,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: IncreaseUserLevelAndExp :exec
UPDATE "User"
SET
	"level" = "level" + 1,
	"exp" = @exp,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: GetInvitedPlayersByReferralCode :many
SELECT
	u.id,
	u.username,
	u.email,
	u.role,
	u."createdAt"
FROM
	"User" u
WHERE
	u."invitedBy" = @invitedBy;

-- name: AddUserBetonPoint :exec
UPDATE "User"
SET
	"betonPoint" = "betonPoint" + @BP,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: SubUserBetonPoint :exec
UPDATE "User"
SET
	"betonPoint" = "betonPoint" - @BP,
	"updatedAt" = now()
WHERE
	id = $1;
