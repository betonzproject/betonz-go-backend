package admin

import (
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
		user, err := auth.Authenticate(app, w, r, "")
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
	BankName      string `formam:"bankName" validate:"oneof=AGD AYA CB KBZ KBZPAY OK_DOLLAR WAVE_PAY YOMA"`
	AccountName   string `formam:"accountName" validate:"required"`
	AccountNumber string `formam:"accountNumber" validate:"number"`
}

type DeleteBankForm struct {
	Id string `formam:"id" validate:"uuid4"`
}

func PostBanks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r, "")
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageSystemBanks) != nil {
			return
		}

		if r.URL.Query().Has("/create") {
			var createBankForm CreateBankForm
			if formutils.ParseDecodeValidate(app, w, r, &createBankForm) != nil {
				return
			}

			app.DB.CreateSystemBank(r.Context(), db.CreateSystemBankParams{
				Name:          db.BankName(createBankForm.BankName),
				AccountName:   createBankForm.AccountName,
				AccountNumber: createBankForm.AccountNumber,
			})
			w.WriteHeader(http.StatusCreated)
			return
		} else if r.URL.Query().Has("/delete") {
			var deleteBankForm DeleteBankForm
			if formutils.ParseDecodeValidate(app, w, r, &deleteBankForm) != nil {
				return
			}

			bankId, _ := utils.ParseUUID(deleteBankForm.Id)
			app.DB.DeleteSystemBankById(r.Context(), bankId)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

type PatchBankForm struct {
	Id            string `formam:"id" validate:"uuid4"`
	AccountName   string `formam:"accountName"`
	AccountNumber string `formam:"accountNumber"`
	Enabled       string `formam:"enabled" validate:"omitempty,oneof=on"`
}

func PatchBanks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r, "")
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
		app.DB.UpdateSystemBank(r.Context(), db.UpdateSystemBankParams{
			ID:            bankId,
			AccountName:   pgtype.Text{String: patchBankForm.AccountName, Valid: patchBankForm.AccountName != ""},
			AccountNumber: pgtype.Text{String: patchBankForm.AccountNumber, Valid: patchBankForm.AccountNumber != ""},
			Disabled:      patchBankForm.Enabled != "on",
		})

		w.WriteHeader(http.StatusOK)
	}
}
