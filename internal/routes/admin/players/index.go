package players

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetPlayers(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManagePlayers) != nil {
			return
		}

		searchParam := r.URL.Query().Get("search")
		fromParam := r.URL.Query().Get("from")
		toParam := r.URL.Query().Get("to")
		statusParam := r.URL.Query().Get("status")

		from, _ := timeutils.ParseDate(fromParam)
		to, err := timeutils.ParseDate(toParam)
		if err != nil {
			to = timeutils.EndOfToday()
		}

		var statuses []db.UserStatus
		if statusParam != "" {
			statuses = []db.UserStatus{db.UserStatus(statusParam)}
		}

		players, err := app.DB.GetUsers(r.Context(), db.GetUsersParams{
			Search:   pgtype.Text{String: searchParam, Valid: searchParam != ""},
			Statuses: statuses,
			FromDate: pgtype.Timestamptz{Time: from, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: to, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get players: " + err.Error())
		}

		jsonutils.Write(w, players, http.StatusOK)
	}
}

type ManageUserForm struct {
	Reason string `form:"reason"`
	UserId string `form:"userId" validate:"uuid4"`
	Status string `form:"status" validate:"oneof=NORMAL RESTRICTED"`
}

func PostPlayers(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManagePlayers) != nil {
			return
		}

		var manageUserForm ManageUserForm
		if formutils.ParseDecodeValidate(app, w, r, &manageUserForm) != nil {
			return
		}

		userToManageId, _ := utils.ParseUUID(manageUserForm.UserId)
		userToManage, err := app.DB.GetUserById(r.Context(), userToManageId)

		if userToManage.Role == db.RoleSUPERADMIN || userToManage.Role == db.RoleADMIN && !acl.IsAuthorized(user.Role, acl.ManageAdmins) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		tx, err := app.Pool.Begin(r.Context())
		if err != nil {
			log.Panicln("Can't start transaction: " + err.Error())
		}
		defer tx.Rollback(r.Context())
		qtx := app.DB.WithTx(tx)

		err = qtx.UpdateUserStatus(r.Context(), db.UpdateUserStatusParams{
			ID:     userToManageId,
			Status: db.UserStatus(manageUserForm.Status),
		})
		if err != nil {
			log.Panicln("Can't update user status: " + err.Error())
		}

		err = qtx.CreateEvent(r.Context(), db.CreateEventParams{
			UserId: user.ID,
			Type:   db.EventTypeCHANGEUSERSTATUS,
			Result: db.EventResultSUCCESS,
			Data: map[string]any{
				"userId": manageUserForm.UserId,
				"status": manageUserForm.Status,
				"reason": manageUserForm.Reason,
			},
		})
		if err != nil {
			log.Panicln("Can't create event: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
