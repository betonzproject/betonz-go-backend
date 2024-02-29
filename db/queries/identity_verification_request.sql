-- name: GetIdentityVerificationRequests :many
SELECT
	ivr.*,
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
		@statuses::"IdentityVerificationStatus"[] IS NULL
		OR ivr.status = ANY (@statuses)
	)
	AND (
		sqlc.narg('search')::TEXT IS NULL
		OR u.username ILIKE '%' || @search || '%'
		OR u2.username ILIKE '%' || @search || '%'
		OR ivr.id::TEXT ILIKE '%' || @search || '%'
		OR ivr."nricName" ILIKE '%' || @search || '%'
		OR ivr.nric ILIKE '%' || @search || '%'
		OR ivr.remarks ILIKE '%' || @search || '%'
	)
	AND ivr."createdAt" >= sqlc.arg('fromDate')
	AND ivr."createdAt" <= sqlc.arg('toDate')
ORDER BY
	ivr.id DESC;

-- name: GetIdentityVerificationRequestById :one
SELECT * FROM "IdentityVerificationRequest" WHERE id = $1;

-- name: GetLatestIdentityVerificationRequestByUserId :one
SELECT * FROM "IdentityVerificationRequest" WHERE "userId" = $1 AND status <> 'REJECTED' ORDER BY "createdAt" DESC LIMIT 1;

-- name: GetPendingIdentityVerificationRequestCount :one
SELECT COUNT(*) FROM "IdentityVerificationRequest" WHERE status = 'PENDING';

-- name: CreateIdentityVerificationRequest :exec
INSERT INTO
	"IdentityVerificationRequest" ("userId", "nricName", nric, dob, "nricFront", "nricBack", "holderFace", status, "updatedAt")
VALUES
	($1, $2, $3, $4, '', '', '', 'INCOMPLETE', now())
RETURNING
	*;

-- name: UpdateIdentityVerificationRequestById :exec
UPDATE
	"IdentityVerificationRequest"
SET
	"modifiedById" = COALESCE(sqlc.narg('modifiedById'), "modifiedById"),
	"nricName" = COALESCE(sqlc.narg('nricName'), "nricName"),
	nric = COALESCE(sqlc.narg('nric'), nric),
	dob = COALESCE(@dob, dob),
	"nricFront" = COALESCE(sqlc.narg('nricFront'), "nricFront"),
	"nricBack" = COALESCE(sqlc.narg('nricBack'), "nricBack"),
	"holderFace" = COALESCE(sqlc.narg('holderFace'), "holderFace"),
	status = COALESCE(sqlc.narg('status'), status),
	remarks = COALESCE(sqlc.narg('remarks'), remarks),
	"updatedAt" = now()
WHERE
	id = $1;
