package admin

import (
	"log"
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/formutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/transactionutils"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetBanks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ToggleSystemBanks) != nil {
			return
		}

		banks, err := app.DB.GetSystemBanks(r.Context())
		jsonutils.Write(w, banks, http.StatusOK)
	}
}

type CreateBankForm struct {
	BankName      string `form:"bankName" validate:"oneof=AGD AYA CB KBZ KBZPAY OK_DOLLAR WAVE_PAY YOMA"`
	AccountName   string `form:"accountName" validate:"required"`
	AccountNumber string `form:"accountNumber" validate:"number"`
}

type DeleteBankForm struct {
	Id string `form:"id" validate:"uuid4"`
}

func PostBanks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageSystemBanks) != nil {
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		if r.URL.Query().Has("/create") {
			var createBankForm CreateBankForm
			if formutils.ParseDecodeValidateMultipart(app, w, r, &createBankForm) != nil {
				return
			}

			_, err = app.DB.GetBankByBankNameAndNumber(r.Context(), db.GetBankByBankNameAndNumberParams{
				Name:          db.BankName(createBankForm.BankName),
				AccountNumber: createBankForm.AccountNumber,
			})
			if err == nil {
				http.Error(w, "bank.alreadyExist.message", http.StatusBadRequest)
				return
			}

			bank, err := qtx.CreateSystemBank(r.Context(), db.CreateSystemBankParams{
				Name:          db.BankName(createBankForm.BankName),
				AccountName:   createBankForm.AccountName,
				AccountNumber: createBankForm.AccountNumber,
			})
			if err != nil {
				log.Panicln("Can't create system bank: " + err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypeSYSTEMBANKADD, db.EventResultSUCCESS, "", map[string]any{
				"bankId": utils.EncodeUUID(bank.ID.Bytes),
				"bank":   string(bank.Name) + " " + string(bank.AccountName) + " " + string(bank.AccountNumber),
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			tx.Commit(r.Context())

			w.WriteHeader(http.StatusCreated)
			return
		} else if r.URL.Query().Has("/delete") {
			var deleteBankForm DeleteBankForm
			if formutils.ParseDecodeValidateMultipart(app, w, r, &deleteBankForm) != nil {
				return
			}

			bankId, _ := utils.ParseUUID(deleteBankForm.Id)
			bank, err := qtx.DeleteSystemBankById(r.Context(), bankId)
			if err != nil {
				log.Panicln("Can't delete system bank: " + err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypeSYSTEMBANKDELETE, db.EventResultSUCCESS, "", map[string]any{
				"bankId": utils.EncodeUUID(bank.ID.Bytes),
				"bank":   string(bank.Name) + " " + string(bank.AccountName) + " " + string(bank.AccountNumber),
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			tx.Commit(r.Context())

			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

type PatchBankForm struct {
	Id            string `form:"id" validate:"uuid4"`
	AccountName   string `form:"accountName"`
	AccountNumber string `form:"accountNumber"`
	Enabled       string `form:"enabled" validate:"omitempty,oneof=on"`
}

func PatchBanks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ToggleSystemBanks) != nil {
			return
		}

		var patchBankForm PatchBankForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &patchBankForm) != nil {
			return
		}

		if (patchBankForm.AccountName != "" || patchBankForm.AccountNumber != "") &&
			acl.Authorize(app, w, r, user.Role, acl.ManageSystemBanks) != nil {
			return
		}

		bankId, _ := utils.ParseUUID(patchBankForm.Id)
		oldBank, err := app.DB.GetSystemBankById(r.Context(), bankId)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		_, err = app.DB.GetBankByBankNameAndNumber(r.Context(), db.GetBankByBankNameAndNumberParams{
			Name:          oldBank.Name,
			AccountNumber: patchBankForm.AccountNumber,
		})
		if err == nil {
			http.Error(w, "bank.alreadyExist.message", http.StatusBadRequest)
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		bank, err := qtx.UpdateSystemBank(r.Context(), db.UpdateSystemBankParams{
			ID:            bankId,
			AccountName:   pgtype.Text{String: patchBankForm.AccountName, Valid: patchBankForm.AccountName != ""},
			AccountNumber: pgtype.Text{String: patchBankForm.AccountNumber, Valid: patchBankForm.AccountNumber != ""},
			Disabled:      patchBankForm.Enabled != "on",
		})
		if err != nil {
			log.Panicln("Can't update system bank: " + err.Error())
		}

		err = utils.LogEvent(qtx, r, user.ID, db.EventTypeSYSTEMBANKUPDATE, db.EventResultSUCCESS, "", map[string]any{
			"bankId":   bankId,
			"old":      string(oldBank.Name) + " " + string(oldBank.AccountName) + " " + string(oldBank.AccountNumber),
			"new":      string(bank.Name) + " " + string(bank.AccountName) + " " + string(bank.AccountNumber),
			"disabled": bank.Disabled,
		})
		if err != nil {
			log.Panicln("Can't log event: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
