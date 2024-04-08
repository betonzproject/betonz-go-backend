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
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
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
		if formutils.ParseDecodeValidateMultipart(app, w, r, &deleteBankForm) != nil {
			return
		}
		bankId, _ := utils.ParseUUID(deleteBankForm.Id)

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		bank, err := qtx.DeleteBankById(r.Context(), bankId)
		if err != nil {
			log.Panicln("Can't delete bank: " + err.Error())
		}

		err = utils.LogEvent(qtx, r, user.ID, db.EventTypeBANKDELETE, db.EventResultSUCCESS, "", map[string]any{
			"bankId": utils.EncodeUUID(bankId.Bytes),
			"bank":   string(bank.Name) + " " + string(bank.AccountName) + " " + string(bank.AccountNumber),
		})
		if err != nil {
			log.Panicln("Can't log event: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
