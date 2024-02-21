package profile

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateForm struct {
	DisplayName string `form:"displayName" validate:"max=30"`
	Email       string `form:"email" validate:"required,email" key:"user.email"`
	CountryCode string `form:"countryCode" validate:"omitempty,number"`
	PhoneNumber string `form:"phoneNumber" validate:"omitempty,number,max=14"`
}

func PostProfile(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if r.URL.Query().Has("/updateProfile") {
			if acl.Authorize(app, w, r, user.Role, acl.UpdateProfile) != nil {
				return
			}

			var updateForm UpdateForm
			phone := ""
			if formutils.ParseDecodeValidate(app, w, r, &updateForm) != nil {
				return
			}
			if updateForm.PhoneNumber != "" {
				phone = "+" + updateForm.CountryCode + updateForm.PhoneNumber
			}

			updateEvent := make(map[string]any)
			if updateForm.DisplayName != user.DisplayName.String {
				updateEvent["displayName"] = updateForm.DisplayName
			}
			if updateForm.Email != user.Email {
				updateEvent["email"] = updateForm.Email
			}
			if phone != user.PhoneNumber.String {
				updateEvent["phoneNumber"] = phone
			}

			tx, qtx := transactionutils.Begin(app, r.Context())
			defer tx.Rollback(r.Context())

			err = qtx.UpdateUser(r.Context(), db.UpdateUserParams{
				ID:          user.ID,
				DisplayName: pgtype.Text{String: updateForm.DisplayName, Valid: updateForm.DisplayName != ""},
				Email:       updateForm.Email,
				PhoneNumber: pgtype.Text{String: phone, Valid: phone != ""},
			})
			if err != nil {
				log.Panicln("Can't update user: " + err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypePROFILEUPDATE, db.EventResultSUCCESS, "", updateEvent)
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			tx.Commit(r.Context())

			w.WriteHeader(http.StatusOK)
			return
		} else if r.URL.Query().Has("/restoreWallet") {
			tx, qtx := transactionutils.Begin(app, r.Context())
			defer tx.Rollback(r.Context())

			errors := restoreWallet(qtx, r.Context(), user)

			err := utils.LogEvent(qtx, r, user.ID, db.EventTypeRESTOREWALLET, db.EventResultSUCCESS, "", map[string]any{
				"errors": errors,
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

func restoreWallet(q *db.Queries, ctx context.Context, user db.User) []string {
	turnoverTargets, err := q.GetTurnoverTargetsByUserId(ctx, user.ID)
	if err != nil {
		log.Panicln("Can't get turnover targets: " + err.Error())
	}

	var refId string
	if os.Getenv("ENVIRONMENT") == "development" {
		refId = "(DEV) TRANSFER"
	} else {
		refId = "TRANSFER"
	}

	var wg sync.WaitGroup

	sum := numericutils.Zero
	errors := make([]string, 0, len(product.AllProducts))
	var sumMutex sync.Mutex
	for _, p := range product.AllProducts {
		if slices.ContainsFunc(turnoverTargets, func(tt db.GetTurnoverTargetsByUserIdRow) bool {
			p2 := product.Product(int(tt.ProductCode.Int32))
			return product.SharesSameWallet(p, p2)
		}) {
			continue
		}

		wg.Add(1)
		go func(p product.Product) {
			defer wg.Done()
			balance, err := product.GetUserBalance(user.EtgUsername, p)
			if err != nil {
				sumMutex.Lock()
				defer sumMutex.Unlock()
				errStr := fmt.Sprintf("Can't get balance of %s (%d) for %s: %s", p, p, user.EtgUsername, err)
				trimmed := strings.Split(errStr, "\n")[0]
				errors = append(errors, trimmed)
				log.Println(errStr)
				return
			}

			if numericutils.IsPositive(balance) {
				err := product.Withdraw(refId, user.EtgUsername, p, balance)

				sumMutex.Lock()
				defer sumMutex.Unlock()
				if err != nil {
					errStr := fmt.Sprintf("Can't transfer from %s (%d) to Main Wallet (-1) for %s: %s", p, p, user.Username, err)
					trimmed := strings.Split(errStr, "\n")[0]
					errors = append(errors, trimmed)
					log.Println(errStr)
				} else {
					sum = numericutils.Add(sum, balance)
				}
			}
		}(p)
	}

	wg.Wait()
	err = q.DepositUserMainWallet(ctx, db.DepositUserMainWalletParams{
		ID:     user.ID,
		Amount: sum,
	})
	if err != nil {
		log.Panicln("Can't restore wallet: " + err.Error())
	}

	return errors
}
