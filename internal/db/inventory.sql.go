// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: inventory.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addItemToInventory = `-- name: AddItemToInventory :exec
INSERT INTO
	"Inventory" ("userId", "item", "count")
VALUES
	($1, $2, $3)
ON CONFLICT ("userId", item) DO
UPDATE
SET
	"count" = "Inventory"."count" + $3
`

type AddItemToInventoryParams struct {
	UserId pgtype.UUID       `json:"userId"`
	Item   InventoryItemType `json:"item"`
	Count  pgtype.Numeric    `json:"count"`
}

func (q *Queries) AddItemToInventory(ctx context.Context, arg AddItemToInventoryParams) error {
	_, err := q.db.Exec(ctx, addItemToInventory, arg.UserId, arg.Item, arg.Count)
	return err
}

const getInventoryByUserId = `-- name: GetInventoryByUserId :many
SELECT
	"userId",
	"item",
	"count"
FROM
	"Inventory"
WHERE
	"userId" = $1
`

type GetInventoryByUserIdRow struct {
	UserId pgtype.UUID       `json:"userId"`
	Item   InventoryItemType `json:"item"`
	Count  pgtype.Numeric    `json:"count"`
}

func (q *Queries) GetInventoryByUserId(ctx context.Context, userid pgtype.UUID) ([]GetInventoryByUserIdRow, error) {
	rows, err := q.db.Query(ctx, getInventoryByUserId, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetInventoryByUserIdRow{}
	for rows.Next() {
		var i GetInventoryByUserIdRow
		if err := rows.Scan(&i.UserId, &i.Item, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
