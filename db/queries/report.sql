-- name: GetDailyReport :many
SELECT * FROM "Report" WHERE "createdAt" >= sqlc.arg('fromDate') AND "createdAt" <= sqlc.arg('toDate');
