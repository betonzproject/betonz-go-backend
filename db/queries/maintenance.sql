-- name: GetMaintenanceList :many
SELECT
	*
FROM
	"Maintenance"
WHERE
	("maintenancePeriod").upper >= CURRENT_TIMESTAMP - interval '1 month'
ORDER BY
	("maintenancePeriod").upper DESC;

-- name: CreateMaintenanceItem :exec
INSERT INTO
	"Maintenance" ("productCode", "maintenancePeriod", "gmtOffsetSecs", "updatedAt")
VALUES
	($1, $2, $3, now());

-- name: DeleteMaintenanceItem :exec
DELETE FROM "Maintenance"
WHERE
	id = $1;

-- name: UpdateMaintenanceItem :exec
UPDATE "Maintenance"
SET
	"maintenancePeriod" = $2,
	"gmtOffsetSecs" = $3,
	"updatedAt" = now()
WHERE
	id = $1;

-- name: GetMaintenanceProductCodes :many
SELECT
	"productCode"
FROM
	"Maintenance"
WHERE
	now() <@ "maintenancePeriod";