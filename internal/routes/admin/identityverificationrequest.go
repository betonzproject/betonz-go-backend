package admin

import (
	"log"
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils/formutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/timeutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/transactionutils"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetIdentityVerificationRequest(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageIdentityVerificationRequests) != nil {
			return
		}

		searchParam := r.URL.Query().Get("search")
		fromParam := r.URL.Query().Get("from")
		toParam := r.URL.Query().Get("to")
		statusParam := r.URL.Query().Get("status")

		from, err := timeutils.ParseDatetime(fromParam)
		if err != nil {
			from = timeutils.StartOfToday().AddDate(0, 0, -6)
		}
		to, err := timeutils.ParseDatetime(toParam)
		if err != nil {
			to = timeutils.EndOfToday()
		}

		var statuses []db.IdentityVerificationStatus
		if statusParam != "" {
			statuses = []db.IdentityVerificationStatus{db.IdentityVerificationStatus(statusParam)}
		}

		requests, err := app.DB.GetIdentityVerificationRequests(r.Context(), db.GetIdentityVerificationRequestsParams{
			Search:   pgtype.Text{String: searchParam, Valid: searchParam != ""},
			Statuses: statuses,
			FromDate: pgtype.Timestamptz{Time: from, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: to, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get identity verification requests: " + err.Error())
		}

		jsonutils.Write(w, requests, http.StatusOK)
	}
}

type IdentityVerificationRequestForm struct {
	RequestId int32  `form:"requestId" validate:"required"`
	Action    string `form:"action" validate:"required,oneof=approve reject"`
	Remarks   string `form:"remarks"`
}

func PostIdentityVerificationRequest(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageIdentityVerificationRequests) != nil {
			return
		}

		var identityVerificationRequestForm IdentityVerificationRequestForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &identityVerificationRequestForm) != nil {
			return
		}

		// Check if identity verification request exists and is pending (or approved if superadmin)
		ivr, err := app.DB.GetIdentityVerificationRequestById(r.Context(), identityVerificationRequestForm.RequestId)
		if err != nil ||
			!(ivr.Status == db.IdentityVerificationStatusPENDING ||
				ivr.Status == db.IdentityVerificationStatusVERIFIED && acl.IsAuthorized(user.Role, acl.OverruleApprovedIdentityVerificationRequests)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var status db.IdentityVerificationStatus
		if identityVerificationRequestForm.Action == "approve" {
			status = db.IdentityVerificationStatusVERIFIED
		} else {
			status = db.IdentityVerificationStatusREJECTED
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		err = qtx.UpdateIdentityVerificationRequestById(r.Context(), db.UpdateIdentityVerificationRequestByIdParams{
			ID:           ivr.ID,
			ModifiedById: user.ID,
			Status:       db.NullIdentityVerificationStatus{IdentityVerificationStatus: status, Valid: true},
			Remarks:      pgtype.Text{String: identityVerificationRequestForm.Remarks, Valid: identityVerificationRequestForm.Remarks != ""},
		})
		if err != nil {
			log.Panicln("Can't update identity verification: " + err.Error())
		}

		err = qtx.UpdateUserDob(r.Context(), db.UpdateUserDobParams{Dob: ivr.Dob, ID: ivr.UserId})
		if err != nil {
			log.Panicln("Can't update user dob: " + err.Error())
		}

		err = qtx.CreateNotification(r.Context(), db.CreateNotificationParams{
			UserId: ivr.UserId,
			Type:   db.NotificationTypeIDENTITYVERIFICATION,
			Variables: map[string]any{
				"id":     ivr.ID,
				"time":   ivr.CreatedAt.Time,
				"action": identityVerificationRequestForm.Action,
			},
		})
		if err != nil {
			log.Panicln("Can't create notification: " + err.Error())
		}
		app.EventServer.Notify(ivr.UserId, "notification")

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
