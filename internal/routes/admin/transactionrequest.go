package admin

import (
	"log"
	"math/big"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetTransactionRequest(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageTransactionRequests) != nil {
			return
		}

		searchParam := r.URL.Query().Get("search")
		fromParam := r.URL.Query().Get("from")
		toParam := r.URL.Query().Get("to")
		transactionTypeParam := r.URL.Query().Get("transactionType")
		statusParam := r.URL.Query().Get("status")

		from, err := timeutils.ParseDate(fromParam)
		if err != nil {
			from = timeutils.StartOfToday().AddDate(0, 0, -6)
		}
		to, err := timeutils.ParseDate(toParam)
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

		requests, err := app.DB.GetTransactionRequests(r.Context(), db.GetTransactionRequestsParams{
			Search:   pgtype.Text{String: searchParam, Valid: searchParam != ""},
			Types:    types,
			Statuses: statuses,
			FromDate: pgtype.Timestamptz{Time: from, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: to, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get requests: " + err.Error())
		}

		jsonutils.Write(w, requests, http.StatusOK)
	}
}

type TransactionRequestForm struct {
	RequestId   int32  `form:"requestId" validate:"required"`
	Action      string `form:"action" validate:"required,oneof=approve decline"`
	Fees        int64  `form:"fees"`
	Remarks     string `form:"remarks"`
	ReceiptData string `form:"receiptData"`
}

func PostTransactionRequest(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageTransactionRequests) != nil {
			return
		}

		var transactionRequestForm TransactionRequestForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &transactionRequestForm) != nil {
			return
		}

		// Check if transaction request exists and is pending
		tr, err := app.DB.GetTransactionRequestById(r.Context(), transactionRequestForm.RequestId)
		if err != nil || tr.Status != db.TransactionStatusPENDING {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if initiator exists and is not restricted
		initiator, err := app.DB.GetExtendedUserById(r.Context(), tr.UserId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if initiator.Status == db.UserStatusRESTRICTED {
			http.Error(w, "user.accountIsRestricted.message", http.StatusBadRequest)
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		if transactionRequestForm.Action == "approve" {
			if tr.Type == db.TransactionTypeDEPOSIT {
				// TODO: Check tr.depositToWallet
				err = qtx.DepositUserMainWallet(r.Context(), db.DepositUserMainWalletParams{
					ID:     initiator.ID,
					Amount: numericutils.Add(tr.Amount, tr.Bonus),
				})
				if err != nil {
					log.Panicln("Can't update user main wallet: " + err.Error())
				}

				err = qtx.UpdateTransactionRequest(r.Context(), db.UpdateTransactionRequestParams{
					ID:           transactionRequestForm.RequestId,
					ModifiedById: user.ID,
					Status:       db.TransactionStatusAPPROVED,
					Remarks:      pgtype.Text{String: transactionRequestForm.Remarks, Valid: transactionRequestForm.Remarks != ""},
				})
				if err != nil {
					log.Panicln("Can't update transaction request: " + err.Error())
				}
			} else {
				// Withdraw
				if numericutils.Cmp(initiator.MainWallet, tr.Amount) < 0 {
					http.Error(w, "transactionRequest.insufficientBalanceInMainWallet.message", http.StatusBadRequest)
					return
				}

				if transactionRequestForm.ReceiptData == "" {
					http.Error(w, "required.message", http.StatusBadRequest)
					return
				}

				err = qtx.WithdrawUserMainWallet(r.Context(), db.WithdrawUserMainWalletParams{
					ID:     initiator.ID,
					Amount: tr.Amount,
				})
				if err != nil {
					log.Panicln("Can't update user main wallet: " + err.Error())
				}

				err = qtx.UpdateTransactionRequest(r.Context(), db.UpdateTransactionRequestParams{
					ID:               transactionRequestForm.RequestId,
					ModifiedById:     user.ID,
					Status:           db.TransactionStatusAPPROVED,
					ReceiptPath:      pgtype.Text{String: transactionRequestForm.ReceiptData, Valid: true},
					WithdrawBankFees: pgtype.Numeric{Int: big.NewInt(transactionRequestForm.Fees), Valid: true},
					Remarks:          pgtype.Text{String: transactionRequestForm.Remarks, Valid: transactionRequestForm.Remarks != ""},
				})
				if err != nil {
					log.Panicln("Can't update transaction request: " + err.Error())
				}
			}
		} else {
			// Decline transaction request
			err = qtx.UpdateTransactionRequest(r.Context(), db.UpdateTransactionRequestParams{
				ID:           transactionRequestForm.RequestId,
				ModifiedById: user.ID,
				Status:       db.TransactionStatusDECLINED,
				Remarks:      pgtype.Text{String: transactionRequestForm.Remarks, Valid: transactionRequestForm.Remarks != ""},
			})
			if err != nil {
				log.Panicln("Can't update transaction request: " + err.Error())
			}
		}

		err = qtx.CreateNotification(r.Context(), db.CreateNotificationParams{
			UserId: tr.UserId,
			Type:   db.NotificationTypeTRANSACTION,
			Variables: map[string]any{
				"id":              transactionRequestForm.RequestId,
				"transactionType": tr.Type,
				"amount":          tr.Amount,
				"action":          transactionRequestForm.Action,
			},
		})
		if err != nil {
			log.Panicln("Can't create notification: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
