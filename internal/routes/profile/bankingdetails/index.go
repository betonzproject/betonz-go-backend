package bankingdetails

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
)

func GetBanks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageOwnBankingDetails) != nil {
			return
		}

		banks, err := app.DB.GetBanksByUserId(r.Context(), user.ID)

		jsonutils.Write(w, banks, http.StatusOK)
	}
}

type DeleteBankForm struct {
	Id string `form:"id" validate:"uuid4"`
}

func DeleteBank(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageOwnBankingDetails) != nil {
			return
		}

		var deleteBankForm DeleteBankForm
		if formutils.ParseDecodeValidate(app, w, r, &deleteBankForm) != nil {
			return
		}
		bankId, _ := utils.ParseUUID(deleteBankForm.Id)

		tx, err := app.Pool.Begin(r.Context())
		if err != nil {
			log.Panicln("Can't start transaction: " + err.Error())
		}
		defer tx.Rollback(r.Context())
		qtx := app.DB.WithTx(tx)

		err = qtx.DeleteBankById(r.Context(), bankId)
		if err != nil {
			log.Panicln("Can't delete bank: " + err.Error())
		}

		err = utils.LogEvent(qtx, r, user.ID, db.EventTypeBANKDELETE, db.EventResultSUCCESS, "", map[string]any{
			"bankId": utils.EncodeUUID(bankId.Bytes),
		})
		if err != nil {
			log.Panicln("Can't log event: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
