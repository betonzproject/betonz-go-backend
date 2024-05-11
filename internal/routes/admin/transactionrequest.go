package admin

import (
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/promotion"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type TransactionRequest struct {
	ID                           int32                `json:"id"`
	UserId                       pgtype.UUID          `json:"userId"`
	ModifiedById                 pgtype.UUID          `json:"modifiedById"`
	BankName                     db.NullBankName      `json:"bankName"`
	BankAccountName              pgtype.Text          `json:"bankAccountName"`
	BankAccountNumber            pgtype.Text          `json:"bankAccountNumber"`
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
	TransactionNo                string               `json:"transactionNo"`
	ModifiedByUsername           pgtype.Text          `json:"modifiedByUsername"`
	ModifiedByRole               db.NullRole          `json:"modifiedByRole"`
}

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

		from, err := timeutils.ParseDatetime(fromParam)
		if err != nil {
			from = timeutils.StartOfToday().AddDate(0, 0, -6)
		}
		to, err := timeutils.ParseDatetime(toParam)
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

		jsonutils.Write(w, sliceutils.Map(requests, func(r db.GetTransactionRequestsRow) TransactionRequest {
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
				TransactionNo:                r.TransactionNo.String,
				ModifiedByUsername:           r.ModifiedByUsername,
				ModifiedByRole:               r.ModifiedByRole,
			}
		}), http.StatusOK)
	}
}

type TransactionRequestForm struct {
	RequestId int32  `form:"requestId" validate:"required"`
	Action    string `form:"action" validate:"required,oneof=approve decline"`
	Fees      int64  `form:"fees"`
	Remarks   string `form:"remarks"`
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
				if tr.Promotion.Valid && tr.DepositToWallet.Valid && tr.DepositToWallet.Int32 != int32(product.MainWallet) {
					var refId string
					if os.Getenv("ENVIRONMENT") == "development" {
						refId = "(DEV) TRANSFER"
					} else {
						refId = "TRANSFER"
					}

					p := product.Product(tr.DepositToWallet.Int32)
					err := product.Deposit(refId, user.EtgUsername, p, numericutils.Add(tr.Amount, tr.Bonus))
					if err != nil {
						log.Printf("Can't deposit to %s (%d) for %s: %s", p, p, user.Username, err)
						http.Error(w, "Deposit failed", http.StatusServiceUnavailable)
						return
					}

					// Promotion turnover
					err = qtx.CreateTurnoverTarget(r.Context(), db.CreateTurnoverTargetParams{
						TransactionRequestId: tr.ID,
						Target:               promotion.CalculateTurnoverTarget(numericutils.Add(tr.Amount, tr.Bonus), tr.Promotion.PromotionType),
					})
					if err != nil {
						log.Panicln("Can't create turnover: " + err.Error())
					}
				} else {
					err := qtx.DepositUserMainWallet(r.Context(), db.DepositUserMainWalletParams{
						ID:     initiator.ID,
						Amount: numericutils.Add(tr.Amount, tr.Bonus),
					})
					if err != nil {
						log.Panicln("Can't update user main wallet: " + err.Error())
					}
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

				depositAmount, _ := tr.Amount.Int64Value()
				if depositAmount.Int64 >= 100000 {
					err = qtx.DepositUserMainWallet(r.Context(), db.DepositUserMainWalletParams{
						ID:     tr.UserId,
						Amount: pgtype.Numeric{Int: big.NewInt(depositAmount.Int64 * 5 / 100), Valid: true},
					})
					if err != nil {
						log.Panicln("Error depositing user's wallet: ", err.Error())
					}
					invitor, err := qtx.GetPlayerByReferralCode(r.Context(), initiator.InvitedBy)
					if err == nil {
						err = qtx.DepositUserMainWallet(r.Context(), db.DepositUserMainWalletParams{
							ID:     invitor.ID,
							Amount: pgtype.Numeric{Int: big.NewInt(depositAmount.Int64 * 5 / 100), Valid: true},
						})
						if err != nil {
							log.Panicln("Error depositing invitor's wallet: ", err.Error())
						}
					}
				}
			} else {
				err = qtx.UpdateTransactionRequest(r.Context(), db.UpdateTransactionRequestParams{
					ID:               transactionRequestForm.RequestId,
					ModifiedById:     user.ID,
					Status:           db.TransactionStatusAPPROVED,
					WithdrawBankFees: pgtype.Numeric{Int: big.NewInt(transactionRequestForm.Fees), Valid: true},
					Remarks:          pgtype.Text{String: transactionRequestForm.Remarks, Valid: transactionRequestForm.Remarks != ""},
				})
				if err != nil {
					log.Panicln("Can't update transaction request: " + err.Error())
				}
			}
		} else {
			if tr.Type == db.TransactionTypeWITHDRAW {
				err = qtx.DepositUserMainWallet(r.Context(), db.DepositUserMainWalletParams{
					ID:     tr.UserId,
					Amount: tr.Amount,
				})
				if err != nil {
					log.Panicln("Error depositing user's MainWallet: ", err.Error())
				}
			}

			// Decline transaction request
			err := qtx.UpdateTransactionRequest(r.Context(), db.UpdateTransactionRequestParams{
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
				"time":            tr.CreatedAt.Time,
				"transactionType": tr.Type,
				"amount":          tr.Amount,
				"action":          transactionRequestForm.Action,
			},
		})
		if err != nil {
			log.Panicln("Can't create notification: " + err.Error())
		}
		app.EventServer.Notify(tr.UserId, "notification")

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
