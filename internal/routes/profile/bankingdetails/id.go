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
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetBankById(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r, "")
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageOwnBankingDetails) != nil {
			return
		}

		bankIdParam := chi.URLParam(r, "bankId")
		bankId, err := utils.ParseUUID(bankIdParam)
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		bank, err := app.DB.GetBankById(r.Context(), bankId)
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		jsonutils.Write(w, bank, http.StatusOK)
	}
}

func PatchBankById(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r, "")
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageOwnBankingDetails) != nil {
			return
		}

		var patchBankForm BankForm
		if formutils.ParseDecodeValidate(app, w, r, &patchBankForm) != nil {
			return
		}

		// Strip spaces
		patchBankForm.AccountNumber = strings.ReplaceAll(patchBankForm.AccountNumber, " ", "")

		bankIdParam := chi.URLParam(r, "bankId")
		bankId, err := utils.ParseUUID(bankIdParam)
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		bank, err := app.DB.GetBankById(r.Context(), bankId)
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		tx, err := app.Pool.Begin(r.Context())
		if err != nil {
			log.Panicln("Can't start transaction: " + err.Error())
		}
		defer tx.Rollback(r.Context())
		qtx := app.DB.WithTx(tx)

		err = qtx.UpdateBank(r.Context(), db.UpdateBankParams{
			ID:            bankId,
			Name:          db.BankName(patchBankForm.BankName),
			AccountName:   pgtype.Text{String: patchBankForm.AccountName, Valid: patchBankForm.AccountName != ""},
			AccountNumber: pgtype.Text{String: patchBankForm.AccountNumber, Valid: patchBankForm.AccountNumber != ""},
		})
		if err != nil {
			log.Panicln("Can't update bank: " + err.Error())
		}

		err = qtx.CreateEvent(r.Context(), db.CreateEventParams{
			SourceIp: pgtype.Text{String: r.RemoteAddr, Valid: true},
			UserId:   user.ID,
			Type:     db.EventTypeBANKUPDATE,
			Result:   db.EventResultSUCCESS,
			Data: map[string]string{
				"bankId": bankIdParam,
				"old":    string(bank.Name) + " " + string(bank.AccountName) + " " + string(bank.AccountNumber),
				"new":    string(patchBankForm.BankName) + " " + string(patchBankForm.AccountName) + " " + string(patchBankForm.AccountNumber),
			},
		})
		if err != nil {
			log.Panicln("Can't create event: " + err.Error())
		}

		tx.Commit(r.Context())

		http.Redirect(w, r, "/profile/banking-details", http.StatusFound)
	}
}
