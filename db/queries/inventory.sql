-- name: AddItemToInventory :exec
INSERT INTO
	"Inventory" ("userId", "item", "count")
VALUES
	($1, $2, $3)
ON CONFLICT ("userId", item) DO
UPDATE
SET
	"count" = "Inventory"."count" + $3;

-- name: GetInventoryByUserId :many
SELECT
	"userId",
	"item",
	"count"
FROM
	"Inventory"
WHERE
	"userId" = $1;