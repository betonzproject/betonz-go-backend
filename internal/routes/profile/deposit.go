package profile

import (
	"crypto/sha1"
	"log"
	"math/big"
	"net/http"
	"slices"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type DepositResponse struct {
	Banks            []db.Bank   `json:"banks"`
	LastUsedBankId   pgtype.UUID `json:"lastUsedBankId"`
	ReceivingBank    *db.Bank    `json:"receivingBank"`
	HasRecentDeposit bool        `json:"hasRecentDeposit"`
}

func GetDeposit(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.DepositToOwnWallet) != nil {
			return
		}

		banks, err := app.DB.GetBanksByUserId(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Can't get banks: " + err.Error())
		}

		// Choose a receiving bank depending on the depositor bank
		depositBankIdParam := r.URL.Query().Get("depositorBankId")
		depositBankId, _ := utils.ParseUUID(depositBankIdParam)
		var depositorBank *db.Bank
		if depositBankId.Valid {
			i := slices.IndexFunc(banks, func(bank db.Bank) bool { return bank.ID.Bytes == depositBankId.Bytes })
			if i != -1 {
				depositorBank = &banks[i]
			}
		} else if user.LastUsedBankId.Valid {
			i := slices.IndexFunc(banks, func(bank db.Bank) bool { return bank.ID.Bytes == user.LastUsedBankId.Bytes })
			if i != -1 {
				depositorBank = &banks[i]
			}
		} else if len(banks) > 0 {
			depositorBank = &banks[0]
		}

		var receivingBank *db.Bank
		if depositorBank != nil {
			systemBanks, err := app.DB.GetSystemBanksByBankName(r.Context(), depositorBank.Name)
			if err != nil {
				log.Panicln("Can't get system banks: " + err.Error())
			}

			if len(systemBanks) > 0 {
				hash := sha1.New()
				h, _ := hash.Write(app.Scs.GetBytes(r.Context(), "sessionId"))

				receivingBank = &systemBanks[h%len(systemBanks)]
			}
		}

		hasRecentDeposit, _ := app.DB.HasRecentDepositRequestsByUserId(r.Context(), user.ID)

		jsonutils.Write(w, DepositResponse{
			Banks:            banks,
			LastUsedBankId:   user.LastUsedBankId,
			ReceivingBank:    receivingBank,
			HasRecentDeposit: hasRecentDeposit,
		}, http.StatusOK)
	}
}

type DepositForm struct {
	DepositorBankId string `form:"depositorBankId" validate:"uuid4"`
	ReceivingBankId string `form:"receivingBankId" validate:"uuid4"`
	DepositAmount   int64  `form:"depositAmount" validate:"min=10000,max=20000000" key:"deposit.amount"`
	ReceiptData     string `form:"receiptData"`
}

func PostDeposit(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.DepositToOwnWallet) != nil {
			return
		}

		var depositForm DepositForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &depositForm) != nil {
			return
		}

		// Check user status
		if user.Status != db.UserStatusNORMAL {
			http.Error(w, "user.accountIsRestricted.message", http.StatusForbidden)
			return
		}

		// Check recent deposits
		hasRecentDeposit, _ := app.DB.HasRecentDepositRequestsByUserId(r.Context(), user.ID)
		if hasRecentDeposit {
			http.Error(w, "deposit.alreadySubmitted.message", http.StatusBadRequest)
			return
		}

		// Validate banks
		depositorBankId, _ := utils.ParseUUID(depositForm.DepositorBankId)
		depositorBank, err := app.DB.GetBankById(r.Context(), depositorBankId)
		if err != nil {
			http.Error(w, "deposit.depositorBankInvalid.message", http.StatusBadRequest)
			return
		}

		receivingBankId, _ := utils.ParseUUID(depositForm.ReceivingBankId)
		receivingBank, err := app.DB.GetSystemBankById(r.Context(), receivingBankId)
		if err != nil || depositorBank.Name != receivingBank.Name {
			http.Error(w, "deposit.receivingBankInvalid.message", http.StatusBadRequest)
			return
		}

		tx, err := app.Pool.Begin(r.Context())
		if err != nil {
			log.Panicln("Can't start transaction: " + err.Error())
		}
		defer tx.Rollback(r.Context())
		qtx := app.DB.WithTx(tx)

		err = qtx.UpdateUserLastUsedBank(r.Context(), db.UpdateUserLastUsedBankParams{
			ID:             user.ID,
			LastUsedBankId: depositorBank.ID,
		})
		if err != nil {
			log.Panicln("Can't update last used bank: " + err.Error())
		}

		err = qtx.CreateTransactionRequest(r.Context(), db.CreateTransactionRequestParams{
			UserId:                       user.ID,
			BankName:                     depositorBank.Name,
			BankAccountName:              depositorBank.AccountName,
			BankAccountNumber:            depositorBank.AccountNumber,
			BeneficiaryBankAccountName:   pgtype.Text{String: receivingBank.AccountName, Valid: true},
			BeneficiaryBankAccountNumber: pgtype.Text{String: receivingBank.AccountNumber, Valid: true},
			Amount:                       pgtype.Numeric{Int: big.NewInt(depositForm.DepositAmount), Valid: true},
			Bonus:                        pgtype.Numeric{Int: big.NewInt(0), Valid: true},
			Type:                         db.TransactionTypeDEPOSIT,
			ReceiptPath:                  pgtype.Text{String: depositForm.ReceiptData, Valid: true},
			Status:                       db.TransactionStatusPENDING,
		})
		if err != nil {
			log.Panicln("Can't create deposit request: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
