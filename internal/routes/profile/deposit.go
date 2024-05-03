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
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/promotion"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type TurnoverTargetInfo struct {
	ProductName   string         `json:"productName"`
	TurnoverSoFar pgtype.Numeric `json:"turnoverSoFar"`
	Target        pgtype.Numeric `json:"target"`
}

type DepositResponse struct {
	Products                  map[product.Product]string `json:"products"`
	Banks                     []db.Bank                  `json:"banks"`
	LastUsedBankId            pgtype.UUID                `json:"lastUsedBankId"`
	ReceivingBank             *db.Bank                   `json:"receivingBank"`
	HasRecentDeposit          bool                       `json:"hasRecentDeposit"`
	EligiblePromotions        []db.PromotionType         `json:"eligiblePromotions"`
	FivePercentBonusRemaining pgtype.Numeric             `json:"fivePercentBonusRemaining"`
	TenPercentBonusRemaining  pgtype.Numeric             `json:"tenPercentBonusRemaining"`
	TurnoverTargets           []TurnoverTargetInfo       `json:"turnoverTargets"`
	ProductsUnderMaintenance  []string                   `json:"productsUnderMaintenance"`
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

		productNames := make(map[product.Product]string)
		for _, p := range product.AllProducts {
			productNames[p] = p.String()
		}

		banks, err := app.DB.GetBanksByUserId(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Can't get banks: " + err.Error())
		}

		// Choose a receiving bank depending on the depositor bank
		bankName := r.URL.Query().Get("bankName")

		var receivingBank *db.Bank
		if bankName != "" {
			systemBanks, err := app.DB.GetSystemBanksByBankName(r.Context(), db.BankName(bankName))
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

		promotions, fivePercentBonusRemaining, tenPercentBonusRemaining := promotion.GetEligiblePromotions(app.DB, r.Context(), user.ID)

		turnoverTargets, err := app.DB.GetTurnoverTargetsByUserId(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Can't get turnover targets: " + err.Error())
		}
		turnoverTargetInfos := sliceutils.Map(turnoverTargets, func(tt db.GetTurnoverTargetsByUserIdRow) TurnoverTargetInfo {
			return TurnoverTargetInfo{
				ProductName:   product.Product(int(tt.ProductCode.Int32)).String(),
				Target:        tt.Target,
				TurnoverSoFar: tt.TurnoverSoFar,
			}
		})

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		jsonutils.Write(w, DepositResponse{
			Products:                  productNames,
			Banks:                     banks,
			LastUsedBankId:            user.LastUsedBankId,
			ReceivingBank:             receivingBank,
			HasRecentDeposit:          hasRecentDeposit,
			EligiblePromotions:        promotions,
			FivePercentBonusRemaining: fivePercentBonusRemaining,
			TenPercentBonusRemaining:  tenPercentBonusRemaining,
			TurnoverTargets:           turnoverTargetInfos,
			ProductsUnderMaintenance: sliceutils.Map(productsUnderMaintenance, func(prodInt int32) string {
				return product.Product(prodInt).String()
			}),
		}, http.StatusOK)
	}
}

type DepositForm struct {
	BankName        string           `form:"bankName" validate:"oneof=AGD AYA CB KBZ KBZPAY OK_DOLLAR WAVE_PAY YOMA"`
	ReceivingBankId string           `form:"receivingBankId" validate:"uuid4"`
	DepositAmount   int64            `form:"depositAmount" validate:"min=10000,max=20000000" key:"deposit.amount"`
	AccountNumber   string           `form:"accountNumber" validate:"required"`
	AccountName     string           `form:"accountName" validate:"required"`
	Promotion       db.PromotionType `form:"promotion" validate:"omitempty,oneof=INACTIVE_BONUS FIVE_PERCENT_UNLIMITED_BONUS TEN_PERCENT_UNLIMITED_BONUS" key:"deposit.promotion"`
	DepositTo       product.Product  `form:"depositTo" validate:"product"`
	ReceiptData     string           `form:"receiptData"`
	TransactionNo   string           `form:"transactionNo"`
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

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		if slices.Contains(productsUnderMaintenance, int32(depositForm.DepositTo)) {
			http.Error(w, "transfer.productUnderMaintenance.message", http.StatusNotAcceptable)
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
			http.Error(w, "deposit.alreadySubmitted.message", http.StatusTooManyRequests)
			return
		}

		receivingBankId, _ := utils.ParseUUID(depositForm.ReceivingBankId)
		receivingBank, err := app.DB.GetSystemBankById(r.Context(), receivingBankId)
		if err != nil || db.BankName(depositForm.BankName) != receivingBank.Name || receivingBank.Disabled {
			http.Error(w, "deposit.receivingBankInvalid.message", http.StatusBadRequest)
			return
		}

		// Validate promotions
		eligiblePromotions, fivePercentBonusRemaining, tenPercentBonusRemaining := promotion.GetEligiblePromotions(app.DB, r.Context(), user.ID)
		if depositForm.Promotion != "" {
			// User must be verified to apply for promotions
			request, _ := app.DB.GetLatestIdentityVerificationRequestByUserId(r.Context(), user.ID)
			if request.Status != db.IdentityVerificationStatusVERIFIED {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if depositForm.DepositTo == product.MainWallet {
				http.Error(w, "deposit.depositToInvalid.message", http.StatusBadRequest)
				return
			}

			if !slices.Contains(eligiblePromotions, depositForm.Promotion) {
				http.Error(w, "deposit.promotionInvalid.message", http.StatusBadRequest)
				return
			}

			hasTurnoverTarget, err := app.DB.HasTurnoverTargetByProductAndUserId(r.Context(), db.HasTurnoverTargetByProductAndUserIdParams{
				UserId:      user.ID,
				ProductCode: pgtype.Int4{Int32: int32(depositForm.DepositTo), Valid: true},
			})
			if err != nil {
				log.Panicln("Can't get turnover target by product: " + err.Error())
			}
			if hasTurnoverTarget {
				http.Error(w, "transfer.unmetTurnoverTarget.message", http.StatusForbidden)
				return
			}
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		amount := pgtype.Numeric{Int: big.NewInt(depositForm.DepositAmount), Valid: true}
		bonus := numericutils.Zero
		if depositForm.Promotion != "" {
			bonus = promotion.CalculateBonus(amount, depositForm.Promotion)
			if depositForm.Promotion == db.PromotionTypeFIVEPERCENTUNLIMITEDBONUS {
				bonus = numericutils.Min(bonus, fivePercentBonusRemaining)
			} else if depositForm.Promotion == db.PromotionTypeTENPERCENTUNLIMITEDBONUS {
				bonus = numericutils.Min(bonus, tenPercentBonusRemaining)
			}
		}

		err = qtx.CreateTransactionRequest(r.Context(), db.CreateTransactionRequestParams{
			UserId: user.ID,
			BankName: db.NullBankName{
				BankName: db.BankName(depositForm.BankName),
				Valid:    true,
			},
			BankAccountName:              pgtype.Text{},
			BankAccountNumber:            pgtype.Text{String: depositForm.AccountNumber, Valid: true},
			BeneficiaryBankAccountName:   pgtype.Text{String: receivingBank.AccountName, Valid: true},
			BeneficiaryBankAccountNumber: pgtype.Text{String: receivingBank.AccountNumber, Valid: true},
			Amount:                       amount,
			Bonus:                        bonus,
			DepositToWallet:              pgtype.Int4{Int32: int32(depositForm.DepositTo), Valid: depositForm.DepositTo != product.MainWallet},
			Type:                         db.TransactionTypeDEPOSIT,
			Promotion:                    db.NullPromotionType{PromotionType: depositForm.Promotion, Valid: depositForm.Promotion != ""},
			ReceiptPath:                  pgtype.Text{String: depositForm.ReceiptData, Valid: true},
			Status:                       db.TransactionStatusPENDING,
			TransactionNo:                pgtype.Text{String: depositForm.TransactionNo, Valid: depositForm.TransactionNo != ""},
		})
		if err != nil {
			log.Panicln("Can't create deposit request: " + err.Error())
		}

		_, err = qtx.GetBankByBankNameAndNumber(r.Context(), db.GetBankByBankNameAndNumberParams{
			AccountNumber: depositForm.AccountNumber,
			Name:          db.BankName(depositForm.BankName),
		})
		if err != nil {
			_, err = qtx.CreateBank(r.Context(), db.CreateBankParams{
				UserId:        user.ID,
				Name:          db.BankName(depositForm.BankName),
				AccountName:   depositForm.AccountName,
				AccountNumber: depositForm.AccountNumber,
			})
			if err != nil {
				log.Panicln("Error creating bank account: ", err.Error())
			}
		}

		app.EventServer.NotifyAdmins("request")

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
