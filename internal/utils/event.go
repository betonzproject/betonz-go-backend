package utils

import (
	"net/http"

	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func LogEvent(q *db.Queries, r *http.Request, userId pgtype.UUID, eventType db.EventType, result db.EventResult, reason string, data map[string]any) error {
	return q.CreateEvent(r.Context(), db.CreateEventParams{
		SourceIp:    pgtype.Text{String: r.RemoteAddr, Valid: r.RemoteAddr != ""},
		UserId:      userId,
		Type:        eventType,
		Result:      result,
		Reason:      pgtype.Text{String: reason, Valid: reason != ""},
		Data:        data,
		HttpRequest: ParseRequest(r),
	})
}
