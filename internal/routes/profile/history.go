package profile

import (
	"log"
	"net/http"
	"time"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type TransactionRequest struct {
	ID                           int32                `json:"id"`
	UserId                       pgtype.UUID          `json:"userId"`
	ModifiedById                 pgtype.UUID          `json:"modifiedById"`
	BankName                     db.BankName          `json:"bankName"`
	BankAccountName              string               `json:"bankAccountName"`
	BankAccountNumber            string               `json:"bankAccountNumber"`
	BeneficiaryBankAccountName   pgtype.Text          `json:"beneficiaryBankAccountName"`
	BeneficiaryBankAccountNumber pgtype.Text          `json:"beneficiaryBankAccountNumber"`
	Amount                       pgtype.Numeric       `json:"amount"`
	Type                         db.TransactionType   `json:"type"`
	ReceiptPath                  pgtype.Text          `json:"receiptPath"`
	Status                       db.TransactionStatus `json:"status"`
	Remarks                      pgtype.Text          `json:"remarks"`
	CreatedAt                    pgtype.Timestamptz   `json:"createdAt"`
	UpdatedAt                    pgtype.Timestamptz   `json:"updatedAt"`
	Bonus                        pgtype.Numeric       `json:"bonus"`
	WithdrawBankFees             pgtype.Numeric       `json:"withdrawBankFees"`
	DepositToWalletName          pgtype.Text          `json:"depositToWalletName"`
	Promotion                    db.NullPromotionType `json:"promotion"`
	Username                     string               `json:"username"`
	Role                         db.Role              `json:"role"`
	ModifiedByUsername           pgtype.Text          `json:"modifiedByUsername"`
	ModifiedByRole               db.NullRole          `json:"modifiedByRole"`
}

func GetHistory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewOwnTransactionHistory) != nil {
			return
		}

		dateRangeParam := r.URL.Query().Get("dateRange")
		transactionTypeParam := r.URL.Query().Get("transactionType")
		statusParam := r.URL.Query().Get("status")

		var from time.Time
		var to time.Time
		from, to, err = timeutils.ParseDateRange(dateRangeParam)
		if err != nil {
			to = timeutils.EndOfToday()
		}

		var types []db.TransactionType
		if transactionTypeParam != "" {
			types = []db.TransactionType{db.TransactionType(transactionTypeParam)}
		}

		var statuses []db.TransactionStatus
		if statusParam != "" {
			statuses = []db.TransactionStatus{db.TransactionStatus(statusParam)}
		}

		requests, err := app.DB.GetTransactionRequestsByUserId(r.Context(), db.GetTransactionRequestsByUserIdParams{
			UserId:   user.ID,
			Types:    types,
			Statuses: statuses,
			FromDate: pgtype.Timestamptz{Time: from, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: to, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get transaction request: " + err.Error())
		}

		jsonutils.Write(w, sliceutils.Map(requests, func(r db.GetTransactionRequestsByUserIdRow) TransactionRequest {
			var depositToWalletName pgtype.Text
			if r.DepositToWallet.Valid {
				depositToWalletName.String = product.Product(int(r.DepositToWallet.Int32)).String()
				depositToWalletName.Valid = true
			}
			return TransactionRequest{
				ID:                           r.ID,
				UserId:                       r.UserId,
				ModifiedById:                 r.ModifiedById,
				BankName:                     r.BankName,
				BankAccountName:              r.BankAccountName,
				BankAccountNumber:            r.BankAccountNumber,
				BeneficiaryBankAccountName:   r.BeneficiaryBankAccountName,
				BeneficiaryBankAccountNumber: r.BeneficiaryBankAccountNumber,
				Amount:                       r.Amount,
				Type:                         r.Type,
				ReceiptPath:                  r.ReceiptPath,
				Status:                       r.Status,
				Remarks:                      r.Remarks,
				CreatedAt:                    r.CreatedAt,
				UpdatedAt:                    r.UpdatedAt,
				Bonus:                        r.Bonus,
				WithdrawBankFees:             r.WithdrawBankFees,
				DepositToWalletName:          depositToWalletName,
				Promotion:                    r.Promotion,
				Username:                     r.Username,
				Role:                         r.Role,
				ModifiedByUsername:           r.ModifiedByUsername,
				ModifiedByRole:               r.ModifiedByRole,
			}
		}), http.StatusOK)
	}
}
