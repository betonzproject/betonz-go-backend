package admin

import (
	"net/http"
	"time"

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
		user, err := auth.Authenticate(app, w, r, "")
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewActivityLog) != nil {
			return
		}

		searchParam := r.URL.Query().Get("search")
		dateRangeParam := r.URL.Query().Get("dateRange")
		roleParam := r.URL.Query().Get("role")
		eventResultParam := r.URL.Query().Get("eventResult")
		eventTypeParam := r.URL.Query().Get("eventType")

		var from time.Time
		var to time.Time
		from, to, err = timeutils.ParseDateRange(dateRangeParam)
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
		jsonutils.Write(w, events, http.StatusOK)
	}
}
