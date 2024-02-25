-- name: GetLatestIdentityVerificationRequestByUserId :one
SELECT * FROM "IdentityVerificationRequest" WHERE "userId" = $1 AND status <> 'REJECTED' ORDER BY "createdAt" DESC LIMIT 1;

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
	"nricName" = COALESCE(sqlc.narg('nricName'), "nricName"),
	nric = COALESCE(sqlc.narg('nric'), nric),
	dob = COALESCE(@dob, dob),
	"nricFront" = COALESCE(sqlc.narg('nricFront'), "nricFront"),
	"nricBack" = COALESCE(sqlc.narg('nricBack'), "nricBack"),
	"holderFace" = COALESCE(sqlc.narg('holderFace'), "holderFace"),
	status = COALESCE(sqlc.narg('status'), status),
	"updatedAt" = now()
WHERE
	id = $1;
