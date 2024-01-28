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

func GetTransactionLog(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewTransactionLogs) != nil {
			return
		}

		searchParam := r.URL.Query().Get("search")
		fromParam := r.URL.Query().Get("from")
		toParam := r.URL.Query().Get("to")

		from, _ := timeutils.ParseDate(fromParam)
		to, err := timeutils.ParseDate(toParam)
		if err != nil {
			to = timeutils.EndOfToday()
		}

		requests, err := app.DB.GetTransactionLogs(r.Context(), db.GetTransactionLogsParams{
			Search:   pgtype.Text{String: searchParam, Valid: searchParam != ""},
			FromDate: pgtype.Timestamptz{Time: from, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: to, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get requests: " + err.Error())
		}

		jsonutils.Write(w, requests, http.StatusOK)
	}
}
