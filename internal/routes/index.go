package routes

import (
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Response struct {
	User                    *db.GetUserByIdRow `json:"user"`
	UnreadNotificationCount int64              `json:"unreadNotificationCount"`
	Permissons              []acl.Permission   `json:"permissions"`
}

func GetIndex(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := app.Scs.GetBytes(r.Context(), "userId")
		if len(userId) < 16 {
			jsonutils.Write(w, Response{}, http.StatusOK)
			return
		}

		user, err := app.DB.GetUserById(r.Context(), pgtype.UUID{Bytes: [16]byte(userId), Valid: true})
		if err != nil {
			jsonutils.Write(w, Response{}, http.StatusOK)
			return
		}

		var unreadNotificationCount int64
		if acl.IsAuthorized(user.Role, acl.ViewNotifications) {
			unreadNotificationCount, err = app.DB.GetUnreadNotificationCountByUserId(r.Context(), user.ID)
		}

		jsonutils.Write(w, Response{User: &user, UnreadNotificationCount: unreadNotificationCount, Permissons: acl.Acl[user.Role]}, http.StatusOK)
	}
}
