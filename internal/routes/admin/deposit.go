package admin

import (
	"log"
	"math/big"
	"net/http"
	"os"
	"slices"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/product"
	"github.com/BetOnz-Company/betonz-go/internal/utils/formutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/numericutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/sliceutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/transactionutils"

	"github.com/jackc/pgx/v5/pgtype"
)

type DepositForm struct {
	Username string          `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Product  product.Product `form:"product" validate:"required"`
	Amount   int64           `form:"amount"  validate:"min=0" key:"deposit.amount"`
}

type DepositResponse struct {
	Products                 map[product.Product]string `json:"products"`
	ProductsUnderMaintenance []string                   `json:"productsUnderMaintenance"`
}

func GetDeposit(app *app.App) http.HandlerFunc {
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

		jsonutils.Write(w, DepositResponse{
			Products: productNames,
			ProductsUnderMaintenance: sliceutils.Map(productsUnderMaintenance, func(prodInt int32) string {
				return product.Product(prodInt).String()
			}),
		}, http.StatusOK)
	}
}

func PostDeposit(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.AdminDeposit) != nil {
			return
		}

		var depositForm DepositForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &depositForm) != nil {
			return
		}

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		if slices.Contains(productsUnderMaintenance, int32(depositForm.Product)) {
			http.Error(w, "deposit.productUnderMaintenance.message", http.StatusNotAcceptable)
			return
		}

		depositAmount := pgtype.Numeric{
			Int:   big.NewInt(depositForm.Amount),
			Valid: true,
		}

		// GetUserByUsername
		userToManage, err := app.DB.GetExtendedUserByUsername(r.Context(), db.GetExtendedUserByUsernameParams{
			Username: depositForm.Username,
		})
		if err != nil {
			http.Error(w, "user.notExist.message", http.StatusBadRequest)
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		// Deposit to Main Wallet
		if depositForm.Product == product.MainWallet {
			err = qtx.DepositUserMainWallet(r.Context(), db.DepositUserMainWalletParams{
				ID:     userToManage.ID,
				Amount: depositAmount,
			})
			if err != nil {
				http.Error(w, "transfer.failed.message", http.StatusServiceUnavailable)
				return
			}

		} else {
			// Deposit to specific wallet
			var refId string
			if os.Getenv("ENVIRONMENT") == "development" {
				refId = "(DEV) TRANSFER"
			} else {
				refId = "TRANSFER"
			}

			err = product.Deposit(refId, userToManage.EtgUsername, depositForm.Product, depositAmount)
			if err != nil {
				http.Error(w, "transfer.failed.message", http.StatusBadRequest)
				return
			}
		}

		err = qtx.CreateTransactionRequest(r.Context(), db.CreateTransactionRequestParams{
			UserId:          userToManage.ID,
			Amount:          depositAmount,
			DepositToWallet: pgtype.Int4{Int32: int32(depositForm.Product), Valid: depositForm.Product != product.MainWallet},
			Type:            db.TransactionTypeDEPOSIT,
			ReceiptPath:     pgtype.Text{Valid: true},
			Bonus:           numericutils.Zero,
			Status:          db.TransactionStatusAPPROVED,
			Remarks:         pgtype.Text{String: "Manual Deposit", Valid: true},
			ModifiedById:    user.ID,
		})
		if err != nil {
			log.Panicln("Error creating transaction request: ", err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusCreated)
	}
}
