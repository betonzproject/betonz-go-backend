// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: event.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createEvent = `-- name: CreateEvent :exec
INSERT INTO
	"Event" ("sourceIp", "userId", type, result, reason, data, "httpRequest", "updatedAt")
VALUES
	($1, $2, $3, $4, $5, $6, $7, now())
`

type CreateEventParams struct {
	SourceIp    pgtype.Text    `json:"sourceIp"`
	UserId      pgtype.UUID    `json:"userId"`
	Type        EventType      `json:"type"`
	Result      EventResult    `json:"result"`
	Reason      pgtype.Text    `json:"reason"`
	Data        map[string]any `json:"data"`
	HttpRequest HttpRequest    `json:"httpRequest"`
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) error {
	_, err := q.db.Exec(ctx, createEvent,
		arg.SourceIp,
		arg.UserId,
		arg.Type,
		arg.Result,
		arg.Reason,
		arg.Data,
		arg.HttpRequest,
	)
	return err
}

const getActiveEventTodayByUserId = `-- name: GetActiveEventTodayByUserId :one
SELECT
	id, "sourceIp", "userId", type, result, reason, data, "createdAt", "updatedAt", "httpRequest"
FROM
	"Event" e
WHERE
	"type" = 'ACTIVE'
	AND e."userId" = $1
	AND date_trunc('day', "createdAt" AT TIME ZONE 'Asia/Yangon') = date_trunc('day', now() AT TIME ZONE 'Asia/Yangon')
LIMIT
	1
`

func (q *Queries) GetActiveEventTodayByUserId(ctx context.Context, userid pgtype.UUID) (Event, error) {
	row := q.db.QueryRow(ctx, getActiveEventTodayByUserId, userid)
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

const getActivePlayerCount = `-- name: GetActivePlayerCount :one
SELECT
	COUNT(DISTINCT e."userId")
FROM
	"User" u
	JOIN "Event" e ON u.id = e."userId"
WHERE
	u.role = 'PLAYER'
	AND e.type = 'ACTIVE'
	AND e."updatedAt" >= $1
	AND e."updatedAt" <= $2
`

type GetActivePlayerCountParams struct {
	FromDate pgtype.Timestamptz `json:"fromDate"`
	ToDate   pgtype.Timestamptz `json:"toDate"`
}

func (q *Queries) GetActivePlayerCount(ctx context.Context, arg GetActivePlayerCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, getActivePlayerCount, arg.FromDate, arg.ToDate)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getEvents = `-- name: GetEvents :many
SELECT
	e.id, e."sourceIp", e."userId", e.type, e.result, e.reason, e.data, e."createdAt", e."updatedAt", e."httpRequest",
	u.username,
	u.role
FROM
	"Event" e
	LEFT JOIN "User" u ON e."userId" = u.id
WHERE
	(
		$1::"Role"[] IS NULL
		OR u.role = ANY ($1)
	)
	AND (
		u.role IS NULL
		OR $2::"Role"[] IS NULL
		OR u.role <> ANY ($2)
	)
	AND e."type" NOT IN ('AUTHENTICATION', 'AUTHORIZATION', 'ACTIVE')
	AND (
		$3::"EventType"[] IS NULL
		OR e."type" = ANY ($3)
	)
	AND (
		$4::"EventResult"[] IS NULL
		OR e.result = ANY ($4)
	)
	AND e."createdAt" >= $5
	AND e."createdAt" <= $6
	AND (
		$7::TEXT IS NULL
		OR u.username ILIKE '%' || $7 || '%'
		OR e."sourceIp" ILIKE '%' || $7 || '%'
		OR e.reason ILIKE '%' || $7 || '%'
		OR e.data::TEXT ILIKE '%' || $7 || '%'
	)
ORDER BY
	e.id DESC
`

type GetEventsParams struct {
	Roles        []Role             `json:"roles"`
	ExcludeRoles []Role             `json:"excludeRoles"`
	Types        []EventType        `json:"types"`
	Results      []EventResult      `json:"results"`
	FromDate     pgtype.Timestamptz `json:"fromDate"`
	ToDate       pgtype.Timestamptz `json:"toDate"`
	Search       pgtype.Text        `json:"search"`
}

type GetEventsRow struct {
	ID          int32              `json:"id"`
	SourceIp    pgtype.Text        `json:"sourceIp"`
	UserId      pgtype.UUID        `json:"userId"`
	Type        EventType          `json:"type"`
	Result      EventResult        `json:"result"`
	Reason      pgtype.Text        `json:"reason"`
	Data        map[string]any     `json:"data"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	HttpRequest HttpRequest        `json:"httpRequest"`
	Username    pgtype.Text        `json:"username"`
	Role        NullRole           `json:"role"`
}

func (q *Queries) GetEvents(ctx context.Context, arg GetEventsParams) ([]GetEventsRow, error) {
	rows, err := q.db.Query(ctx, getEvents,
		arg.Roles,
		arg.ExcludeRoles,
		arg.Types,
		arg.Results,
		arg.FromDate,
		arg.ToDate,
		arg.Search,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEventsRow{}
	for rows.Next() {
		var i GetEventsRow
		if err := rows.Scan(
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
			&i.Username,
			&i.Role,
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

const getRestrictionEventsByUserId = `-- name: GetRestrictionEventsByUserId :many
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
	AND data ->> 'userId' = $1::UUID::TEXT
ORDER BY
	e."createdAt" DESC
`

type GetRestrictionEventsByUserIdRow struct {
	ID        int32              `json:"id"`
	UserId    pgtype.UUID        `json:"userId"`
	Username  string             `json:"username"`
	Role      Role               `json:"role"`
	Data      map[string]any     `json:"data"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

func (q *Queries) GetRestrictionEventsByUserId(ctx context.Context, userid pgtype.UUID) ([]GetRestrictionEventsByUserIdRow, error) {
	rows, err := q.db.Query(ctx, getRestrictionEventsByUserId, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetRestrictionEventsByUserIdRow{}
	for rows.Next() {
		var i GetRestrictionEventsByUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.Username,
			&i.Role,
			&i.Data,
			&i.CreatedAt,
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

const updateEvent = `-- name: UpdateEvent :exec
UPDATE "Event" SET "updatedAt" = now() WHERE id = $1
`

func (q *Queries) UpdateEvent(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, updateEvent, id)
	return err
}
