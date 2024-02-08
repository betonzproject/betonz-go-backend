-- name: GetEvents :many
SELECT
	e.*,
	u.username,
	u.role
FROM
	"Event" e
	LEFT JOIN "User" u ON e."userId" = u.id
WHERE
	(
		@roles::"Role"[] IS NULL
		OR u.role = ANY (@roles)
	)
	AND (
		u.role IS NULL
		OR sqlc.arg('excludeRoles')::"Role"[] IS NULL
		OR u.role <> ANY (sqlc.arg('excludeRoles'))
	)
	AND e."type" NOT IN ('AUTHENTICATION', 'AUTHORIZATION', 'ACTIVE')
	AND (
		@types::"EventType"[] IS NULL
		OR e."type" = ANY (@types)
	)
	AND (
		@results::"EventResult"[] IS NULL
		OR e.result = ANY (@results)
	)
	AND e."createdAt" >= sqlc.arg('fromDate')
	AND e."createdAt" <= sqlc.arg('toDate')
	AND (
		sqlc.narg('search')::TEXT IS NULL
		OR u.username ILIKE '%' || @search || '%'
		OR e."sourceIp" ILIKE '%' || @search || '%'
		OR e.reason ILIKE '%' || @search || '%'
		OR e.data::TEXT ILIKE '%' || @search || '%'
	)
ORDER BY
	e.id DESC;

-- name: GetRestrictionEventsByUserId :many
SELECT
	e.id,
	e."userId",
	u.username,
	u.role,
	e.data,
	e."createdAt"
FROM
	"Event" e
	JOIN "User" u ON e."userId" = u.id
WHERE
	type = 'CHANGE_USER_STATUS'
	AND data ->> 'userId' = sqlc.arg('userId')::UUID::TEXT
ORDER BY
	e."createdAt" DESC;

-- name: CreateEvent :exec
INSERT INTO
	"Event" ("sourceIp", "userId", type, result, reason, data, "httpRequest", "updatedAt")
VALUES
	($1, $2, $3, $4, $5, $6, $7, now());
