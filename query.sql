-- name: GetUserById :one
SELECT id, username, role, email, "displayName", "phoneNumber", "mainWallet", dob, "profileImage", "isEmailVerified" FROM "User" WHERE id = $1;

-- name: GetExtendedUserByUsername :one
SELECT * FROM "User" WHERE username = $1 AND (@roles::"Role"[] IS NULL OR role = ANY(@roles));



-- name: GetTransactionRequests :many
SELECT
	tr.*,
	u.username,
	u.role,
	u2.username AS "modifiedByUsername",
	u2.role AS "modifiedByRole"
FROM
	"TransactionRequest" tr
JOIN "User" u ON
	u.id = tr."userId"
LEFT JOIN "User" u2 ON
	u2.id = tr."modifiedById"
WHERE 
	(@types::"TransactionType"[] IS NULL OR tr."type" = ANY(@types))
AND
	(@statuses::"TransactionStatus"[] IS NULL OR tr.status = ANY(@statuses))
AND (
	sqlc.narg('search')::text IS NULL
	OR u.username ILIKE '%' || @search || '%' 
	OR u2.username ILIKE '%' || @search || '%'
	OR tr."bankAccountName" ILIKE '%' || @search || '%'
	OR tr."beneficiaryBankAccountName" ILIKE '%' || @search || '%'
	OR tr.remarks ILIKE '%' || @search || '%'
)
AND
	tr."createdAt" >= sqlc.arg('fromDate')
AND
	tr."createdAt" < sqlc.arg('toDate')
ORDER BY
	tr.id DESC;



-- name: GetUsers :many
SELECT
	*
FROM (
	SELECT
		ROW_NUMBER() OVER (ORDER BY "createdAt") "rowNumber",
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
		u."lastLoginIp",
		u."createdAt",
		st."lastLogin"::timestamptz AS "lastLogin",
		tr1."lastDeposit"::timestamptz AS "lastDeposit",
		tr2."lastWithdraw"::timestamptz AS "lastWithdraw"
	FROM
		"User" u
	LEFT JOIN (
		-- Get last login time
		SELECT "userId", max("createdAt") "lastLogin" FROM "SessionToken" GROUP BY "userId"
	) st ON 
		u.id = st."userId"
	LEFT JOIN (
		-- Get last deposit time
		SELECT
			"userId",
			max("updatedAt") "lastDeposit"
		FROM
			"TransactionRequest"
		WHERE
			"type" = 'DEPOSIT'::"TransactionType"
		AND
			"status" = 'APPROVED'::"TransactionStatus"
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
			"type" = 'WITHDRAW'::"TransactionType"
		AND
			"status" = 'APPROVED'::"TransactionStatus"
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



-- name: UpdateUserStatus :exec
UPDATE "User" SET status = $2, "updatedAt" = now() WHERE id = $1;
