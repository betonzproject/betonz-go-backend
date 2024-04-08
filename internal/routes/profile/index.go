package profile

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/mailutils"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetProfileResponse struct {
	*db.IdentityVerificationStatus `json:"identityVerificationStatus"`
}

func GetProfile(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetUser(app, w, r)
		if err != nil {
			jsonutils.Write(w, GetProfileResponse{}, http.StatusOK)
			return
		}

		request, err := app.DB.GetLatestIdentityVerificationRequestByUserId(r.Context(), user.ID)
		if err != nil {
			jsonutils.Write(w, GetProfileResponse{}, http.StatusOK)
			return
		}

		jsonutils.Write(w, GetProfileResponse{IdentityVerificationStatus: &request.Status}, http.StatusOK)
	}
}

type UpdateProfileForm struct {
	DisplayName string `form:"displayName" validate:"max=30"`
	Email       string `form:"email" validate:"required,email" key:"user.email"`
	CountryCode string `form:"countryCode" validate:"omitempty,number"`
	PhoneNumber string `form:"phoneNumber" validate:"omitempty,number,max=14"`
}

type PostProfileResponse struct {
	ResentVerification bool `json:"resentVerification"`
	ProfileUpdate      bool `json:"profileUpdate"`
}

func PostProfile(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		if r.URL.Query().Has("/update") {
			if acl.Authorize(app, w, r, user.Role, acl.UpdateProfile) != nil {
				return
			}

			var updateProfileForm UpdateProfileForm
			phone := ""
			if formutils.ParseDecodeValidateMultipart(app, w, r, &updateProfileForm) != nil {
				return
			}
			if updateProfileForm.PhoneNumber != "" {
				phone = "+" + updateProfileForm.CountryCode + updateProfileForm.PhoneNumber
			}

			updateEvent := make(map[string]any)
			if updateProfileForm.DisplayName != user.DisplayName.String {
				updateEvent["displayName"] = updateProfileForm.DisplayName
			}
			if user.PendingEmail.Valid && updateProfileForm.Email != user.PendingEmail.String || !user.PendingEmail.Valid && updateProfileForm.Email != user.Email {
				updateEvent["email"] = updateProfileForm.Email

				var templateData struct {
					Subject string
					Body    string
				}

				cookie, err := r.Cookie("i18next")
				var lng string
				if err != nil {
					lng = "en"
				} else {
					lng = cookie.Value
				}

				if lng == "my" {
					templateData = struct {
						Subject string
						Body    string
					}{
						Subject: "လျှိဝှက်နံပါတ်ပြန်လည်သတ်မှတ်",
						Body: `
							<p>မင်္ဂလာပါ ` + user.Username + `ရေ,</p>
							<p>သင်၏ အီးမေးအား ` + updateProfileForm.Email + ` သိုပြောင်းလဲလိုက်သည် ။ 
							သင်၏email တွင် စစ်ဆေးကြည့်ပါ ။ </p>
							<p>သင်မဟုတ်ပါက Customers Service သို ချက်ချင်းဆက်သွယ်ပါ။</p></a>`,
					}
				} else {
					templateData = struct {
						Subject string
						Body    string
					}{
						Subject: "Email Change",
						Body: `
							<p>Hello ` + user.Username + `,</p>
							<p>The email of your account was just changed to ` + updateProfileForm.Email + `. 
							Please check your email to verify this new email address. 
							If you didn't request this change, please contact us immediately.</p>`,
					}
				}

				body, err := utils.ParseTemplate("template.html", templateData)
				if err != nil {
					log.Panicln("Can't parse template: ", err.Error())
				}

				go func() {
					err := mailutils.SendMail(user.Email, body, templateData.Subject)
					if err != nil {
						log.Println("Can't send mail: " + err.Error())
					}
				}()

				requestVerification(qtx, r, user, updateProfileForm.Email)
			}
			if phone != user.PhoneNumber.String {
				updateEvent["phoneNumber"] = phone
			}

			err = qtx.UpdateUser(r.Context(), db.UpdateUserParams{
				ID:           user.ID,
				DisplayName:  pgtype.Text{String: updateProfileForm.DisplayName, Valid: updateProfileForm.DisplayName != ""},
				PendingEmail: pgtype.Text{String: updateProfileForm.Email, Valid: updateProfileForm.Email != ""},
				PhoneNumber:  pgtype.Text{String: phone, Valid: phone != ""},
			})
			if err != nil {
				log.Panicln("Can't update user: " + err.Error())
			}

			err = utils.LogEvent(qtx, r, user.ID, db.EventTypePROFILEUPDATE, db.EventResultSUCCESS, "", updateEvent)
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			tx.Commit(r.Context())

			jsonutils.Write(w, PostProfileResponse{ProfileUpdate: true}, http.StatusOK)
			return
		} else if r.URL.Query().Has("/resendVerification") && user.PendingEmail.Valid {
			requestVerification(qtx, r, user, user.PendingEmail.String)

			tx.Commit(r.Context())

			jsonutils.Write(w, PostProfileResponse{ResentVerification: true}, http.StatusOK)
			return
		} else if r.URL.Query().Has("/restoreWallet") {
			errors := restoreWallet(qtx, r.Context(), user)

			err := utils.LogEvent(qtx, r, user.ID, db.EventTypeRESTOREWALLET, db.EventResultSUCCESS, "", map[string]any{
				"errors": errors,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			tx.Commit(r.Context())

			jsonutils.Write(w, PostProfileResponse{}, http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

func requestVerification(q *db.Queries, r *http.Request, user db.User, newEmail string) {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	token := base64.RawURLEncoding.EncodeToString(randomBytes)

	hash := sha256.New()
	hash.Write([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

	err := q.UpsertVerificationToken(r.Context(), db.UpsertVerificationTokenParams{
		TokenHash: tokenHash,
		UserId:    user.ID,
	})
	if err != nil {
		log.Panicln("Can't create verification token: ", err)
	}

	var templateData struct {
		Subject string
		Body    string
	}

	cookie, err := r.Cookie("i18next")
	var lng string
	if err != nil {
		lng = "en"
	} else {
		lng = cookie.Value
	}
	href := r.Header.Get("Origin") + "/verify-email/" + token

	if lng == "my" {
		templateData = struct {
			Subject string
			Body    string
		}{
			Subject: "အီးမေးအတည်ပြု",
			Body: `
				<p>မင်္ဂလာပါ ` + user.Username + `ရေ,</p>
				<p>သင်၏ အီးမေးအား ` + newEmail + ` သို ချိန်းပြီးပါပြီ။ အောက်ပါလင့်ကို နှိပ်ပီး အတည်ပြုပေးပါ။ လင့်ထဲသို တစ်နာရီအတွင်း သာ ဝင်ရောက်ခွင့်ရှိမည်</p>
				<center style="margin-top: 10px;"><button style="color:white;background:#f3b83d;padding:.5rem .8rem;border-radius:999px;border:none"><a style="color:black;text-decoration:none" href="` + href + "\">Verify Email</a></button></center>",
		}
	} else {
		templateData = struct {
			Subject string
			Body    string
		}{
			Subject: "Verify Email",
			Body: `
				<p>Hello ` + user.Username + `,</p>
				<p>You have requested to verify your email ` + newEmail + `. Click the link below to verify your email. The link will expire in 1 hour.</p>
				<center style="margin-top: 10px;"><button style="color:white;background:#f3b83d;padding:.5rem .8rem;border-radius:999px;border:none"><a style="color:black;text-decoration:none" href="` + href + "\">Verify Email</a></button></center>",
		}
	}

	body, err := utils.ParseTemplate("template.html", templateData)
	if err != nil {
		log.Panicln("Can't parse template : ", err.Error())
	}

	go func() {
		err := mailutils.SendMail(newEmail, body, templateData.Subject)
		if err != nil {
			log.Println("Can't send mail: " + err.Error())
		}
	}()
}

func restoreWallet(q *db.Queries, ctx context.Context, user db.User) []string {
	productsUnderMaintenance, err := q.GetMaintenanceProductCodes(ctx)
	if err != nil {
		log.Panicln("Can't get maintained products: " + err.Error())
	}

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
			if slices.Contains(productsUnderMaintenance, int32(p)) {
				errStr := fmt.Sprintf("Can't get balance of %s (%d) for %s: %s", p, p, user.EtgUsername, err)
				trimmed := strings.Split(errStr, "\n")[0]
				errors = append(errors, trimmed)
				log.Println(errStr)
				return
			}

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
