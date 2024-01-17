-- name: GetNotificationsByUserId :many
SELECT * FROM "Notification" WHERE "userId" = $1 ORDER BY "createdAt" DESC;

-- name: GetUnreadNotificationCountByUserId :one
SELECT count(*) FROM "Notification" WHERE "userId" = $1 AND NOT read;



-- name: MarkNotificationsAsReadByUserId :exec
UPDATE "Notification" SET read = TRUE WHERE "userId" = $1;



-- name: DeleteNotificationById :exec
DELETE FROM "Notification" WHERE id = $1;
