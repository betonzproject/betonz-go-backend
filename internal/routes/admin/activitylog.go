package admin

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetActivityLog(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewActivityLog) != nil {
			return
		}

		searchParam := r.URL.Query().Get("search")
		fromParam := r.URL.Query().Get("from")
		toParam := r.URL.Query().Get("to")
		roleParam := r.URL.Query().Get("role")
		eventResultParam := r.URL.Query().Get("eventResult")
		eventTypeParam := r.URL.Query().Get("eventType")

		from, err := timeutils.ParseDate(fromParam)
		if err != nil {
			from = timeutils.StartOfToday()
		}
		to, err := timeutils.ParseDate(toParam)
		if err != nil {
			to = timeutils.EndOfToday()
		}

		var roles []db.Role
		if roleParam != "" {
			roles = []db.Role{db.Role(roleParam)}
		}

		var excludes []db.Role
		if !acl.IsAuthorized(user.Role, acl.ViewSuperadminActivityLog) {
			excludes = []db.Role{db.RoleSUPERADMIN}
		}

		var eventTypes []db.EventType
		if eventTypeParam != "" {
			eventTypes = []db.EventType{db.EventType(eventTypeParam)}
		}

		var eventResults []db.EventResult
		if eventResultParam != "" {
			eventResults = []db.EventResult{db.EventResult(eventResultParam)}
		}

		events, err := app.DB.GetEvents(r.Context(), db.GetEventsParams{
			Search:       pgtype.Text{String: searchParam, Valid: searchParam != ""},
			Roles:        roles,
			ExcludeRoles: excludes,
			Types:        eventTypes,
			Results:      eventResults,
			FromDate:     pgtype.Timestamptz{Time: from, Valid: true},
			ToDate:       pgtype.Timestamptz{Time: to, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get events: " + err.Error())
		}

		jsonutils.Write(w, events, http.StatusOK)
	}
}
