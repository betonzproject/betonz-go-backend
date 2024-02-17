package profile

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateForm struct {
	DisplayName string `form:"displayName" validate:"max=30"`
	Email       string `form:"email" validate:"required,email" key:"user.email"`
	CountryCode string `form:"countryCode" validate:"omitempty,number"`
	PhoneNumber string `form:"phoneNumber" validate:"omitempty,number,max=14"`
}

func PostProfile(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.UpdateProfile) != nil {
			return
		}

		var updateForm UpdateForm
		phone := ""
		if formutils.ParseDecodeValidate(app, w, r, &updateForm) != nil {
			return
		}
		if updateForm.PhoneNumber != "" {
			phone = "+" + updateForm.CountryCode + updateForm.PhoneNumber
		}

		updateEvent := make(map[string]any)
		if updateForm.DisplayName != user.DisplayName.String {
			updateEvent["displayName"] = updateForm.DisplayName
		}
		if updateForm.Email != user.Email {
			updateEvent["email"] = updateForm.Email
		}
		if phone != user.PhoneNumber.String {
			updateEvent["phoneNumber"] = phone
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		err = qtx.UpdateUser(r.Context(), db.UpdateUserParams{
			ID:          user.ID,
			DisplayName: pgtype.Text{String: updateForm.DisplayName, Valid: updateForm.DisplayName != ""},
			Email:       updateForm.Email,
			PhoneNumber: pgtype.Text{String: phone, Valid: phone != ""},
		})
		if err != nil {
			log.Panicln("Can't update user: " + err.Error())
		}

		err = utils.LogEvent(qtx, r, user.ID, db.EventTypePROFILEUPDATE, db.EventResultSUCCESS, "", updateEvent)
		if err != nil {
			log.Panicln("Can't log event: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
