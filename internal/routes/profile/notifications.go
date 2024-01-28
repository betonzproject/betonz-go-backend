package profile

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
)

func GetNotifications(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewNotifications) != nil {
			return
		}

		notifications, err := app.DB.GetNotificationsByUserId(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Can't get notifications: " + err.Error())
		}

		jsonutils.Write(w, notifications, http.StatusOK)
	}
}

type DeleteNotificationForm struct {
	Id int32 `form:"id" validate:"numeric"`
}

func PostNotifications(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewNotifications) != nil {
			return
		}

		if r.URL.Query().Has("/delete") {
			var deleteNotificationForm DeleteNotificationForm
			if formutils.ParseDecodeValidate(app, w, r, &deleteNotificationForm) != nil {
				return
			}

			err = app.DB.DeleteNotificationById(r.Context(), deleteNotificationForm.Id)
			if err != nil {
				log.Panicln("Can't delete notification: " + err.Error())
			}
		} else {
			err = app.DB.MarkNotificationsAsReadByUserId(r.Context(), user.ID)
			if err != nil {
				log.Panicln("Can't mark notifications as read: " + err.Error())
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
