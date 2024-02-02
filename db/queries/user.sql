-- name: GetUserById :one
SELECT id, username, role, email, "displayName", "phoneNumber", "mainWallet", dob, "profileImage" FROM "User" WHERE id = $1;

-- name: GetExtendedUserById :one
SELECT * FROM "User" WHERE id = $1;

-- name: GetExtendedUserByUsername :one
SELECT * FROM "User" WHERE username = $1 AND (@roles::"Role"[] IS NULL OR role = ANY(@roles));

-- name: GetUsers :many
SELECT
	*
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
		(@statuses::"UserStatus"[] IS NULL OR u.status = ANY(@statuses))
	AND (
		sqlc.narg('search')::text IS NULL
		OR u.username ILIKE '%' || @search || '%' 
		OR u.email ILIKE '%' || @search || '%'
		OR u."displayName" ILIKE '%' || @search || '%'
		OR u."phoneNumber" ILIKE '%' || @search || '%'
		OR u."lastLoginIp" ILIKE '%' || @search || '%'
	)
	AND
		u."createdAt" >= sqlc.arg('fromDate')
	AND
		u."createdAt" < sqlc.arg('toDate')
) q
ORDER BY
	"rowNumber" DESC, "createdAt" DESC;

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
	u.id = $1;

-- name: UpdateUser :exec
UPDATE "User" SET "displayName" = $2, email = $3, "phoneNumber" = $4, "updatedAt" = now() WHERE id = $1;

-- name: UpdateUserProfileImage :exec
UPDATE "User" SET "profileImage" = $2, "updatedAt" = now() WHERE id = $1;

-- name: UpdateUserLastUsedBank :exec
UPDATE "User" SET "lastUsedBankId" = $2, "updatedAt" = now() WHERE id = $1;

-- name: UpdateUserStatus :exec
UPDATE "User" SET status = $2, "updatedAt" = now() WHERE id = $1;
