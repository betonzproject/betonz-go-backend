-- name: GetLastRewardClaimedById :one
SELECT
	*
FROM
	"Event"
WHERE
	type = 'REWARD_CLAIM'
	AND "userId" = $1
ORDER BY
	"createdAt" DESC;

-- name: HasRecentClaimedRewardByUserId :one
SELECT
	EXISTS (
		SELECT
			*
		FROM
			"Event"
		WHERE
			"userId" = $1
			AND type = 'REWARD_CLAIM'
			AND "createdAt" >= now() - INTERVAL '24 hour'
	);
