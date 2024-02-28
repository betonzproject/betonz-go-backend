package admin

import (
	"log"
	"math/big"
	"net/http"
	"os"
	"slices"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type WithdrawForm struct {
	Username     string          `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Product      product.Product `form:"product" validate:"required"`
	Amount       int64           `form:"amount" validate:"min=1" key:"withdraw.amount"`
	Remarks      string          `form:"remarks"`
	OtherRemarks string          `form:"otherRemarks"`
}

type WithdrawResponse struct {
	Products                 map[product.Product]string `json:"products"`
	ProductsUnderMaintenance []string                   `json:"productsUnderMaintenance"`
}

func GetWithdraw(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.AdminDeposit) != nil {
			return
		}

		productNames := make(map[product.Product]string)

		for _, p := range product.AllProducts {
			productNames[p] = p.String()
		}

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		jsonutils.Write(w, WithdrawResponse{
			Products: productNames,
			ProductsUnderMaintenance: sliceutils.Map(productsUnderMaintenance, func(prodInt int32) string {
				return product.Product(prodInt).String()
			}),
		}, http.StatusOK)
	}
}

func PostWithdraw(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.AdminWithdraw) != nil {
			return
		}

		var withdrawForm WithdrawForm
		if formutils.ParseDecodeValidate(app, w, r, &withdrawForm) != nil {
			return
		}

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		if slices.Contains(productsUnderMaintenance, int32(withdrawForm.Product)) {
			http.Error(w, "deposit.productUnderMaintenance.message", http.StatusNotAcceptable)
			return

		}

		// GetUserByUsername
		userToManage, err := app.DB.GetExtendedUserByUsername(r.Context(), db.GetExtendedUserByUsernameParams{
			Username: withdrawForm.Username,
		})
		if err != nil {
			http.Error(w, "user.notExist.message", http.StatusBadRequest)
			return
		}

		withdrawAmount := pgtype.Numeric{
			Int:   big.NewInt(withdrawForm.Amount),
			Valid: true,
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		// Withdraw from Main Wallet
		if withdrawForm.Product == product.MainWallet {
			if numericutils.Cmp(userToManage.MainWallet, withdrawAmount) < 0 {
				http.Error(w, "withdraw.insufficientBalanceInMainWallet.message", http.StatusBadRequest)
				return
			}

			err = qtx.WithdrawUserMainWallet(r.Context(), db.WithdrawUserMainWalletParams{
				ID:     userToManage.ID,
				Amount: withdrawAmount,
			})
			if err != nil {
				http.Error(w, "transfer.failed.message", http.StatusServiceUnavailable)
				return
			}

		} else {
			// withdraw from specific wallet
			var refId string
			if os.Getenv("ENVIRONMENT") == "development" {
				refId = "(DEV) TRANSFER"
			} else {
				refId = "TRANSFER"
			}

			err = product.Withdraw(refId, userToManage.EtgUsername, withdrawForm.Product, withdrawAmount)
			if err != nil {
				http.Error(w, "transfer.failed.message", http.StatusServiceUnavailable)
				return
			}
		}

		var remarks string
		if withdrawForm.Remarks == "Others" {
			remarks = withdrawForm.OtherRemarks
		} else {
			remarks = withdrawForm.Remarks
		}

		err = qtx.CreateTransactionRequest(r.Context(), db.CreateTransactionRequestParams{
			UserId:          userToManage.ID,
			Amount:          withdrawAmount,
			DepositToWallet: pgtype.Int4{Int32: int32(withdrawForm.Product), Valid: withdrawForm.Product != product.MainWallet},
			Type:            db.TransactionTypeWITHDRAW,
			ReceiptPath:     pgtype.Text{Valid: true},
			Bonus:           numericutils.Zero,
			Status:          db.TransactionStatusAPPROVED,
			Remarks:         pgtype.Text{String: remarks, Valid: remarks != ""},
			ModifiedById:    user.ID,
		})
		if err != nil {
			log.Panicln("Error creating transaction request: ", err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
