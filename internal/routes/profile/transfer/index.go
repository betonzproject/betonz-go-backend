package transfer

import (
	"log"
	"math/big"
	"net/http"
	"os"
	"slices"
	"sync"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/product"
	"github.com/BetOnz-Company/betonz-go/internal/utils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/formutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/numericutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/transactionutils"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetTransfer(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.TransferBetweenWallets) != nil {
			return
		}

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		balances := getAllBalances(user.EtgUsername, productsUnderMaintenance)

		jsonutils.Write(w, balances, http.StatusOK)
	}
}

func getAllBalances(etgUsername string, productsUnderMaintenance []int32) map[product.Product]pgtype.Numeric {
	var wg sync.WaitGroup

	balances := make(map[product.Product]pgtype.Numeric)
	for _, p := range product.AllProducts {
		balances[p] = pgtype.Numeric{}
	}
	var balancesMutex sync.Mutex
	for _, p := range product.AllProducts {
		wg.Add(1)
		go func(p product.Product) {
			defer wg.Done()
			if slices.Contains(productsUnderMaintenance, int32(p)) {
				return
			}
			balance, err := product.GetUserBalance(etgUsername, p)
			if err != nil {
				log.Printf("Can't get balance of %s (%d) for %s: %s\n", p, p, etgUsername, err)
				return
			}
			balancesMutex.Lock()
			defer balancesMutex.Unlock()
			balances[p] = balance
		}(p)
	}

	wg.Wait()

	return balances
}

type TransferForm struct {
	FromWallet product.Product `form:"fromWallet" validate:"required"`
	ToWallet   product.Product `form:"toWallet" validate:"required"`
	Amount     int64           `form:"amount" validate:"min=1" key:"transfer.amount"`
}

func PostTransfer(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		var transferForm TransferForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &transferForm) != nil {
			return
		}

		if product.SharesSameWallet(transferForm.FromWallet, transferForm.ToWallet) {
			w.WriteHeader(http.StatusOK)
			return
		}

		data := map[string]any{
			"fromWallet": transferForm.FromWallet,
			"toWallet":   transferForm.ToWallet,
			"amount":     transferForm.Amount,
		}

		productsUnderMaintenance, err := qtx.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintenance product codes: ", err.Error())
		}

		if slices.Contains(productsUnderMaintenance, int32(transferForm.FromWallet)) || slices.Contains(productsUnderMaintenance, int32(transferForm.ToWallet)) {
			http.Error(w, "transfer.productUnderMaintenance.message", http.StatusForbidden)
			return
		}

		// Check turnover target
		if transferForm.FromWallet != product.MainWallet {
			turnoverTargets, err := qtx.GetTurnoverTargetsByUserId(r.Context(), user.ID)
			if err != nil {
				log.Panicln("Can't get turnover targets: " + err.Error())
			}

			if slices.ContainsFunc(turnoverTargets, func(tt db.GetTurnoverTargetsByUserIdRow) bool {
				p := product.Product(int(tt.ProductCode.Int32))
				return product.SharesSameWallet(p, transferForm.FromWallet)
			}) {
				tx.Commit(r.Context())
				http.Error(w, "transfer.unmetTurnoverTarget.message", http.StatusForbidden)
				return
			}
		}

		var fromWalletBalance pgtype.Numeric
		if transferForm.FromWallet == product.MainWallet {
			fromWalletBalance = user.MainWallet
		} else {
			fromWalletBalance, err = product.GetUserBalance(user.EtgUsername, transferForm.FromWallet)
			if err != nil {
				err := utils.LogEvent(app.DB, r, user.ID, db.EventTypeTRANSFERWALLET, db.EventResultFAIL, err.Error(), data)
				if err != nil {
					log.Panicln("Can't log event: " + err.Error())
				}

				log.Printf("Can't get balance of %s (%d) for %s: %s\n", transferForm.FromWallet, transferForm.FromWallet, user.EtgUsername, err)
				http.Error(w, "transfer.failed.message", http.StatusForbidden)
				return
			}
		}

		amount := pgtype.Numeric{Int: big.NewInt(transferForm.Amount), Valid: true}

		if numericutils.Cmp(fromWalletBalance, amount) < 0 {
			err := utils.LogEvent(app.DB, r, user.ID, db.EventTypeTRANSFERWALLET, db.EventResultFAIL, "Insufficient balance", data)
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			http.Error(w, "transfer.insufficientBalance.message", http.StatusForbidden)
			return
		}

		var refId string
		if os.Getenv("ENVIRONMENT") == "development" {
			refId = "(DEV) TRANSFER"
		} else {
			refId = "TRANSFER"
		}

		err = product.Transfer(qtx, r.Context(), refId, user, transferForm.FromWallet, transferForm.ToWallet, amount)
		if err != nil {
			if transferForm.FromWallet != product.MainWallet && transferForm.ToWallet != product.MainWallet {
				// It's possible that as a last resort, the amount gets deposited back to the main wallet,
				// so we need to commit changes done in `product.Transfer()`
				err := utils.LogEvent(qtx, r, user.ID, db.EventTypeTRANSFERWALLET, db.EventResultFAIL, err.Error(), data)
				if err != nil {
					log.Panicln("Can't log event: " + err.Error())
				}
				tx.Commit(r.Context())
			} else {
				err := utils.LogEvent(app.DB, r, user.ID, db.EventTypeTRANSFERWALLET, db.EventResultFAIL, err.Error(), data)
				if err != nil {
					log.Panicln("Can't log event: " + err.Error())
				}
			}

			log.Printf("Can't transfer from %s (%d) to %s (%d) for %s: %s\n", transferForm.FromWallet, transferForm.FromWallet, transferForm.ToWallet, transferForm.ToWallet, user.Username, err)
			http.Error(w, "transfer.failed.message", http.StatusForbidden)
			return
		}

		err = utils.LogEvent(qtx, r, user.ID, db.EventTypeTRANSFERWALLET, db.EventResultSUCCESS, "", data)
		if err != nil {
			log.Panicln("Can't log event: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
