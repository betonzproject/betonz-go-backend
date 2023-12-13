-- name: GetEvents :many
SELECT
	e.*,
	u.username,
	u.role
FROM
	"Event" e
LEFT JOIN "User" u ON
	e."userId" = u.id
WHERE
	(@roles::"Role"[] IS NULL OR u.role = ANY(@roles))
AND
	(u.role IS NULL OR sqlc.arg('excludeRoles')::"Role"[] IS NULL OR u.role <> ANY(sqlc.arg('excludeRoles')))
AND
	e."type" NOT IN ('AUTHENTICATION'::"EventType", 'AUTHORIZATION'::"EventType", 'ACTIVE'::"EventType")
AND
	(@types::"EventType"[] IS NULL OR e."type" = ANY(@types))
AND
	(@results::"EventResult"[] IS NULL OR e.result = ANY(@results))
AND
	e."createdAt" >= sqlc.arg('fromDate')
AND
	e."createdAt" < sqlc.arg('toDate')
AND (
	sqlc.narg('search')::text IS NULL
	OR u.username ILIKE '%' || @search || '%' 
	OR e."sourceIp" ILIKE '%' || @search || '%'
	OR e.reason ILIKE '%' || @search || '%'
	OR e.data::text ILIKE '%' || @search || '%'
)
ORDER BY
	e.id DESC;
