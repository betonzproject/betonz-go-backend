package profile

import (
	"log"
	"math/big"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type WithdrawResponse struct {
	Banks             []db.Bank   `json:"banks"`
	LastUsedBankId    pgtype.UUID `json:"lastUsedBankId"`
	HasRecentWithdraw bool        `json:"hasRecentWithdraw"`
}

func GetWithdraw(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.WithdrawFromOwnWallet) != nil {
			return
		}

		banks, err := app.DB.GetBanksByUserId(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Can't get banks: " + err.Error())
		}

		hasRecentWithdraw, _ := app.DB.HasRecentWithdrawRequestsByUserId(r.Context(), user.ID)

		jsonutils.Write(w, WithdrawResponse{
			Banks:             banks,
			LastUsedBankId:    user.LastUsedBankId,
			HasRecentWithdraw: hasRecentWithdraw,
		}, http.StatusOK)
	}
}

type WithdrawForm struct {
	WithdrawerBankId string `form:"withdrawerBankId" validate:"uuid4"`
	WithdrawAmount   int64  `form:"withdrawAmount" validate:"min=10000,max=20000000" key:"withdraw.amount"`
}

func PostWithdraw(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.WithdrawFromOwnWallet) != nil {
			return
		}

		var withdrawForm WithdrawForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &withdrawForm) != nil {
			return
		}

		// Check user status
		if user.Status != db.UserStatusNORMAL {
			http.Error(w, "user.accountIsRestricted.message", http.StatusForbidden)
			return
		}

		// Check recent withdraws
		hasRecentWithdraw, _ := app.DB.HasRecentWithdrawRequestsByUserId(r.Context(), user.ID)
		if hasRecentWithdraw {
			http.Error(w, "withdraw.alreadySubmitted.message", http.StatusBadRequest)
			return
		}

		// Validate banks
		withdrawerBankId, _ := utils.ParseUUID(withdrawForm.WithdrawerBankId)
		withdrawerBank, err := app.DB.GetBankById(r.Context(), withdrawerBankId)
		if err != nil {
			http.Error(w, "withdraw.withdrawerBankInvalid.message", http.StatusBadRequest)
			return
		}

		// Check balance
		withdrawAmount := pgtype.Numeric{Int: big.NewInt(withdrawForm.WithdrawAmount), Valid: true}
		if numericutils.Cmp(user.MainWallet, withdrawAmount) < 0 {
			http.Error(w, "withdraw.insufficientBalanceInMainWallet.message", http.StatusBadRequest)
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		err = qtx.UpdateUserLastUsedBank(r.Context(), db.UpdateUserLastUsedBankParams{
			ID:             user.ID,
			LastUsedBankId: withdrawerBank.ID,
		})
		if err != nil {
			log.Panicln("Can't update last used bank: " + err.Error())
		}

		err = qtx.CreateTransactionRequest(r.Context(), db.CreateTransactionRequestParams{
			UserId: user.ID,
			BankName: db.NullBankName{
				BankName: withdrawerBank.Name,
				Valid:    true,
			},
			BankAccountName:   pgtype.Text{String: withdrawerBank.AccountName, Valid: true},
			BankAccountNumber: pgtype.Text{String: withdrawerBank.AccountNumber, Valid: true},
			Amount:            withdrawAmount,
			Bonus:             pgtype.Numeric{Int: big.NewInt(0), Valid: true},
			Type:              db.TransactionTypeWITHDRAW,
			Status:            db.TransactionStatusPENDING,
		})
		if err != nil {
			log.Panicln("Can't create withdraw request: " + err.Error())
		}

		app.EventServer.NotifyAdmins("request")

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
