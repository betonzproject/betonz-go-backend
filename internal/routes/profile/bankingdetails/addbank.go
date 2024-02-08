package bankingdetails

import (
	"log"
	"net/http"
	"strings"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
)

type BankForm struct {
	BankName      string `form:"bankName" validate:"oneof=AGD AYA CB KBZ KBZPAY OK_DOLLAR WAVE_PAY YOMA"`
	AccountName   string `form:"accountName" validate:"required"`
	AccountNumber string `form:"accountNumber" validate:"required,accountnumber" key:"bank.accountNumber"`
}

func AddBank(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageOwnBankingDetails) != nil {
			return
		}

		var addBankForm BankForm
		if formutils.ParseDecodeValidate(app, w, r, &addBankForm) != nil {
			return
		}

		// Strip spaces
		addBankForm.AccountNumber = strings.ReplaceAll(addBankForm.AccountNumber, " ", "")

		tx, err := app.Pool.Begin(r.Context())
		if err != nil {
			log.Panicln("Can't start transaction: " + err.Error())
		}
		defer tx.Rollback(r.Context())
		qtx := app.DB.WithTx(tx)

		bank, err := qtx.CreateBank(r.Context(), db.CreateBankParams{
			UserId:        user.ID,
			Name:          db.BankName(addBankForm.BankName),
			AccountName:   addBankForm.AccountName,
			AccountNumber: addBankForm.AccountNumber,
		})
		if err != nil {
			log.Panicln("Can't create bank: " + err.Error())
		}

		err = utils.LogEvent(qtx, r, user.ID, db.EventTypeBANKADD, db.EventResultSUCCESS, "", map[string]any{
			"bankId": utils.EncodeUUID(bank.ID.Bytes),
			"bank":   string(bank.Name) + " " + string(bank.AccountName) + " " + string(bank.AccountNumber),
		})
		if err != nil {
			log.Panicln("Can't create event: " + err.Error())
		}

		tx.Commit(r.Context())

		http.Redirect(w, r, "/profile/banking-details", http.StatusFound)
	}
}
