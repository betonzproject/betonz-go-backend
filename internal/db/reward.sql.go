// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: reward.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getLastRewardClaimedById = `-- name: GetLastRewardClaimedById :one
SELECT
	id, "sourceIp", "userId", type, result, reason, data, "createdAt", "updatedAt", "httpRequest"
FROM
	"Event"
WHERE
	type = 'REWARD_CLAIM'
	AND "userId" = $1
ORDER BY
	"createdAt" DESC
`

func (q *Queries) GetLastRewardClaimedById(ctx context.Context, userid pgtype.UUID) (Event, error) {
	row := q.db.QueryRow(ctx, getLastRewardClaimedById, userid)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.SourceIp,
		&i.UserId,
		&i.Type,
		&i.Result,
		&i.Reason,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HttpRequest,
	)
	return i, err
}

const hasRecentClaimedRewardByUserId = `-- name: HasRecentClaimedRewardByUserId :one
SELECT
	EXISTS (
		SELECT
			id, "sourceIp", "userId", type, result, reason, data, "createdAt", "updatedAt", "httpRequest"
		FROM
			"Event"
		WHERE
			"userId" = $1
			AND type = 'REWARD_CLAIM'
			AND "createdAt" >= now() - INTERVAL '24 hour'
	)
`

func (q *Queries) HasRecentClaimedRewardByUserId(ctx context.Context, userid pgtype.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, hasRecentClaimedRewardByUserId, userid)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
