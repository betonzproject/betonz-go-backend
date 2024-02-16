package admin

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

		tx, err := app.Pool.Begin(r.Context())
		if err != nil {
			log.Panicln("Can't start transaction: " + err.Error())
		}
		defer tx.Rollback(r.Context())
		qtx := app.DB.WithTx(tx)

		if r.URL.Query().Has("/create") {
			var createBankForm CreateBankForm
			if formutils.ParseDecodeValidate(app, w, r, &createBankForm) != nil {
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
			if formutils.ParseDecodeValidate(app, w, r, &deleteBankForm) != nil {
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
		if formutils.ParseDecodeValidate(app, w, r, &patchBankForm) != nil {
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

		tx, err := app.Pool.Begin(r.Context())
		if err != nil {
			log.Panicln("Can't start transaction: " + err.Error())
		}
		defer tx.Rollback(r.Context())
		qtx := app.DB.WithTx(tx)

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
