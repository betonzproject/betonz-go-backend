package profile

import (
	"log"
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

func GetHistory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewOwnTransactionHistory) != nil {
			return
		}

		dateRangeParam := r.URL.Query().Get("dateRange")
		transactionTypeParam := r.URL.Query().Get("transactionType")
		statusParam := r.URL.Query().Get("status")

		var from time.Time
		var to time.Time
		from, to, err = timeutils.ParseDateRange(dateRangeParam)
		if err != nil {
			to = timeutils.EndOfToday()
		}

		var types []db.TransactionType
		if transactionTypeParam != "" {
			types = []db.TransactionType{db.TransactionType(transactionTypeParam)}
		}

		var statuses []db.TransactionStatus
		if statusParam != "" {
			statuses = []db.TransactionStatus{db.TransactionStatus(statusParam)}
		}

		requests, err := app.DB.GetTransactionRequestsByUserId(r.Context(), db.GetTransactionRequestsByUserIdParams{
			UserId:   user.ID,
			Types:    types,
			Statuses: statuses,
			FromDate: pgtype.Timestamptz{Time: from, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: to, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get transaction request: " + err.Error())
		}

		jsonutils.Write(w, requests, http.StatusOK)
	}
}
