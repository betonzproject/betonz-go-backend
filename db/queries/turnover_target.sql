-- name: HasActivePromotionByUserId :one
SELECT
	EXISTS (
		SELECT
			*
		FROM
			"TurnoverTarget"
		WHERE
			"userId" = $1
			AND "promoCode" = $2
			AND (
				sqlc.narg('productCode')::int IS NULL
				OR "productCode" = sqlc.narg('productCode')
			)
	);
