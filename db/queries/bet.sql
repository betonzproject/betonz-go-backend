-- name: GetTurnoverByUserId :many
SELECT
	b."productCode", sum(b.turnover) AS turnover
FROM
	"Bet" b
JOIN
	"User" u USING ("etgUsername")
WHERE 
	u.id = $1 AND b."startTime" >= sqlc.arg('fromDate') AND b."startTime" < sqlc.arg('toDate')
GROUP BY
	b."productCode";
