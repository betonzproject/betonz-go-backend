-- name: GetNotificationsByUserId :many
SELECT * FROM "Notification" WHERE "userId" = $1 ORDER BY "createdAt" DESC;

-- name: GetUnreadNotificationCountByUserId :one
SELECT count(*) FROM "Notification" WHERE "userId" = $1 AND NOT read;

-- name: CreateNotification :exec
INSERT INTO "Notification" ("userId", type, message, variables, "updatedAt") VALUES ($1, $2, $3, $4, now());

-- name: MarkNotificationsAsReadByUserId :exec
UPDATE "Notification" SET read = TRUE WHERE "userId" = $1;

-- name: DeleteNotificationById :exec
DELETE FROM "Notification" WHERE id = $1;