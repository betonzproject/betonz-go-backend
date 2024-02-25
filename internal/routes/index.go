package routes

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
)

type Response struct {
	User                    *db.GetUserByIdRow `json:"user"`
	UnreadNotificationCount int64              `json:"unreadNotificationCount"`
	Permissons              []acl.Permission   `json:"permissions"`
}

func GetIndex(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetUser(app, w, r)
		if err != nil {
			jsonutils.Write(w, Response{}, http.StatusOK)
			return
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

		jsonutils.Write(w, Response{User: &user, UnreadNotificationCount: unreadNotificationCount, Permissons: acl.Acl[user.Role]}, http.StatusOK)
	}
}
