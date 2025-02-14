// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: notification.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createNotification = `-- name: CreateNotification :exec
INSERT INTO "Notification" ("userId", type, message, variables, "updatedAt") VALUES ($1, $2, $3, $4, now())
`

type CreateNotificationParams struct {
	UserId    pgtype.UUID      `json:"userId"`
	Type      NotificationType `json:"type"`
	Message   pgtype.Text      `json:"message"`
	Variables map[string]any   `json:"variables"`
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) error {
	_, err := q.db.Exec(ctx, createNotification,
		arg.UserId,
		arg.Type,
		arg.Message,
		arg.Variables,
	)
	return err
}

const deleteNotificationById = `-- name: DeleteNotificationById :exec
DELETE FROM "Notification" WHERE id = $1
`

func (q *Queries) DeleteNotificationById(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteNotificationById, id)
	return err
}

const getNotificationsByUserId = `-- name: GetNotificationsByUserId :many
SELECT id, "userId", type, message, variables, read, "createdAt", "updatedAt" FROM "Notification" WHERE "userId" = $1 ORDER BY "createdAt" DESC
`

func (q *Queries) GetNotificationsByUserId(ctx context.Context, userid pgtype.UUID) ([]Notification, error) {
	rows, err := q.db.Query(ctx, getNotificationsByUserId, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Notification{}
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.Type,
			&i.Message,
			&i.Variables,
			&i.Read,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getUnreadNotificationCountByUserId = `-- name: GetUnreadNotificationCountByUserId :one
SELECT count(*) FROM "Notification" WHERE "userId" = $1 AND NOT read
`

func (q *Queries) GetUnreadNotificationCountByUserId(ctx context.Context, userid pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getUnreadNotificationCountByUserId, userid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const markNotificationsAsReadByUserId = `-- name: MarkNotificationsAsReadByUserId :exec
UPDATE "Notification" SET read = TRUE WHERE "userId" = $1
`

func (q *Queries) MarkNotificationsAsReadByUserId(ctx context.Context, userid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, markNotificationsAsReadByUserId, userid)
	return err
}
