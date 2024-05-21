package routes

import (
	"log"
	"math/big"
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"

	"github.com/jackc/pgx/v5/pgtype"
)

type Response struct {
	User                    *db.GetUserByIdRow `json:"user"`
	UnreadNotificationCount int64              `json:"unreadNotificationCount"`
	Permissons              []acl.Permission   `json:"permissions"`
	ExpTarget               int64              `json:"expTarget"`
}

func GetIndex(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetUser(app, w, r)
		if err != nil {
			jsonutils.Write(w, Response{}, http.StatusOK)
			return
		}

		userLevel, _ := user.Level.Int64Value()
		userExp, _ := user.Exp.Int64Value()

		nextLevelExp := utils.AllTargets[userLevel.Int64-1]

		canIncreaseLevel := utils.ExpTarget(userExp.Int64) >= nextLevelExp

		if canIncreaseLevel && userLevel.Int64 != 80 {
			err = app.DB.IncreaseUserLevelAndExp(r.Context(), db.IncreaseUserLevelAndExpParams{
				ID:  user.ID,
				Exp: pgtype.Numeric{Int: big.NewInt(int64(utils.ExpTarget(userExp.Int64) - nextLevelExp)), Valid: true},
			})
			if err != nil {
				log.Panicln("Error updating user's level: ", err.Error())
			}
		}

		var unreadNotificationCount int64
		if acl.IsAuthorized(user.Role, acl.ViewNotifications) {
			unreadNotificationCount, err = app.DB.GetUnreadNotificationCountByUserId(r.Context(), user.ID)
			if err != nil {
				log.Panicln("Can't get notification count: ", err)
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

		if canIncreaseLevel && userLevel.Int64 != 80 {
			jsonutils.Write(w, Response{User: &user, UnreadNotificationCount: unreadNotificationCount, Permissons: acl.Acl[user.Role], ExpTarget: int64(utils.ExpTarget(userExp.Int64) - nextLevelExp)}, http.StatusOK)
		} else {
			jsonutils.Write(w, Response{User: &user, UnreadNotificationCount: unreadNotificationCount, Permissons: acl.Acl[user.Role], ExpTarget: int64(nextLevelExp)}, http.StatusOK)
		}
	}
}
