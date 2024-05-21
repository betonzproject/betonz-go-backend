package admin

import (
	"log"
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
)

type Response struct {
	User                                    *db.GetUserByIdRow `json:"user"`
	PendingTransactionRequestCount          int64              `json:"pendingTransactionRequestCount"`
	PendingIdentityVerificationRequestCount int64              `json:"pendingIdentityVerificationRequestCount"`
	Permissons                              []acl.Permission   `json:"permissions"`
}

func GetIndex(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetUser(app, w, r)
		if err != nil {
			jsonutils.Write(w, Response{}, http.StatusOK)
			return
		}

		var pendingTransactionRequestCount int64
		if acl.IsAuthorized(user.Role, acl.ManageTransactionRequests) {
			pendingTransactionRequestCount, err = app.DB.GetPendingTransactionRequestCount(r.Context())
			if err != nil {
				log.Panicln("Can't get pending transaction request count: " + err.Error())
			}
		}

		var pendingIdentityVerificationRequestCount int64
		if acl.IsAuthorized(user.Role, acl.ManageIdentityVerificationRequests) {
			pendingIdentityVerificationRequestCount, err = app.DB.GetPendingIdentityVerificationRequestCount(r.Context())
			if err != nil {
				log.Panicln("Can't get pending identity verification request count: " + err.Error())
			}
		}

		event, err := app.DB.GetActiveEventTodayByUserId(r.Context(), user.ID)
		if err != nil {
			err = utils.LogEvent(app.DB, r, user.ID, db.EventTypeACTIVE, db.EventResultSUCCESS, "", nil)
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}
		} else {
			err = app.DB.UpdateEvent(r.Context(), event.ID)
			if err != nil {
				log.Panicln("Can't update event: " + err.Error())
			}
		}

		jsonutils.Write(w, Response{
			User:                                    &user,
			PendingTransactionRequestCount:          pendingTransactionRequestCount,
			PendingIdentityVerificationRequestCount: pendingIdentityVerificationRequestCount,
			Permissons:                              acl.Acl[user.Role],
		}, http.StatusOK)
	}
}
