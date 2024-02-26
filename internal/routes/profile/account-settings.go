package profile

import (
	"log"
	"net/http"
	"time"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/ratelimiter"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
)

type UpdateUsernameForm struct {
	Username string `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
}

type UpdatePasswordForm struct {
	Password    string `form:"password" validate:"required,min=8,max=512"`
	NewPassword string `form:"newPassword" validate:"required,min=8,max=512"`
}

type AccountSettingsResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var usernameChangeLimitOpts = ratelimiter.LimiterOptions{
	Tokens: 1,
	Window: time.Duration(24 * time.Hour),
}

func PostAccountSettings(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		if r.URL.Query().Has("/username") {
			var updateUsernameForm UpdateUsernameForm
			if formutils.ParseDecodeValidate(app, w, r, &updateUsernameForm) != nil {
				return
			}

			if acl.Authorize(app, w, r, user.Role, acl.UpdateProfile) != nil {
				return
			}

			if user.Username == updateUsernameForm.Username {
				jsonutils.Write(w, AccountSettingsResponse{Type: "username"}, http.StatusOK)
				return
			}

			key := "username_change:" + utils.EncodeUUID(user.ID.Bytes)
			err = app.Limiter.Consume(r.Context(), key, usernameChangeLimitOpts)
			if err == ratelimiter.RateLimited {
				jsonutils.Write(w, AccountSettingsResponse{
					Type:    "username",
					Message: "accountSettings.usernameChangedTooManyTimes.message",
				}, http.StatusTooManyRequests)
				return
			}

			exisitingUser, err := qtx.GetExtendedUserByUsername(r.Context(), db.GetExtendedUserByUsernameParams{
				Username: updateUsernameForm.Username,
			})
			if err == nil && exisitingUser.ID != user.ID {
				err := utils.LogEvent(app.DB, r, user.ID, db.EventTypeUSERNAMECHANGE, db.EventResultFAIL, "Username taken", map[string]any{
					"username": updateUsernameForm.Username,
				})
				if err != nil {
					log.Panicln("Can't log event: ", err.Error())
				}

				jsonutils.Write(w, AccountSettingsResponse{
					Type:    "username",
					Message: "user.username.alreadyTaken.message",
				}, http.StatusForbidden)
				return
			}

			err = qtx.UpdateUsername(r.Context(), db.UpdateUsernameParams{
				ID:       user.ID,
				Username: updateUsernameForm.Username,
			})
			if err != nil {
				log.Panicln("Can't update username: ", err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypeUSERNAMECHANGE, db.EventResultSUCCESS, "", map[string]any{
				"old": user.Username,
				"new": updateUsernameForm.Username,
			})
			if err != nil {
				log.Panicln("Can't log event: ", err.Error())
			}

			tx.Commit(r.Context())

			jsonutils.Write(w, AccountSettingsResponse{
				Type:    "username",
				Message: "accountSettings.usernameChanged.message",
			}, http.StatusOK)
			return
		} else if r.URL.Query().Has("/password") {
			var updatePasswordForm UpdatePasswordForm
			if formutils.ParseDecodeValidate(app, w, r, &updatePasswordForm) != nil {
				return
			}

			passwordMatches, _ := utils.Argon2IDVerify(updatePasswordForm.Password, user.PasswordHash)
			if !passwordMatches {
				err := utils.LogEvent(app.DB, r, user.ID, db.EventTypePASSWORDCHANGE, db.EventResultSUCCESS, "", nil)
				if err != nil {
					log.Panicln("Can't log event: ", err.Error())
				}

				jsonutils.Write(w, AccountSettingsResponse{
					Type:    "password",
					Message: "accountSettings.passwordIncorrect.message",
				}, http.StatusUnauthorized)
				return
			}

			hashedPassword, err := utils.Argon2IDHash(updatePasswordForm.NewPassword)
			if err != nil {
				log.Panicln("Can't hash password: ", err.Error())
			}

			err = qtx.UpdateUserPasswordHash(r.Context(), db.UpdateUserPasswordHashParams{
				ID:           user.ID,
				PasswordHash: hashedPassword,
			})
			if err != nil {
				log.Panicln("Can't update password: ", err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypePASSWORDCHANGE, db.EventResultSUCCESS, "", nil)
			if err != nil {
				log.Panicln("Can't log event: ", err.Error())
			}

			tx.Commit(r.Context())

			jsonutils.Write(w, AccountSettingsResponse{Type: "password"}, http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}
