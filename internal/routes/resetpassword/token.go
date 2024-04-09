package resetpassword

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/mailutils"
	"github.com/doorman2137/betonz-go/internal/utils/ratelimiter"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyTokenResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func GetPasswordResetToken(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")

		hash := sha256.New()
		hash.Write([]byte(token))
		tokenHash := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

		passwordResetToken, err := checkPasswordResetToken(app, r, tokenHash)
		if err != nil {
			if errors.Is(err, ratelimiter.RateLimited) {
				jsonutils.Write(w, VerifyTokenResponse{Message: "tooManyRequests.message"}, http.StatusTooManyRequests)
			} else {
				jsonutils.Write(w, VerifyTokenResponse{Message: "resetPassword.passwordResetLinkInvalid.message"}, http.StatusBadRequest)
			}
			return
		}

		jsonutils.Write(w, VerifyTokenResponse{Username: passwordResetToken.Username}, http.StatusOK)
	}
}

type ResetPasswordForm struct {
	Password string `form:"password" validate:"required,min=8,max=512"`
}

func PostPasswordResetToken(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resetPasswordForm ResetPasswordForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &resetPasswordForm) != nil {
			return
		}

		token := chi.URLParam(r, "token")

		hash := sha256.New()
		hash.Write([]byte(token))
		tokenHash := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

		passwordResetToken, err := checkPasswordResetToken(app, r, tokenHash)
		if err != nil {
			if errors.Is(err, ratelimiter.RateLimited) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
			} else {
				http.Error(w, "Invalid token", http.StatusBadRequest)
			}
			return
		}

		passwordHash, err := utils.Argon2IDHash(resetPasswordForm.Password)
		if err != nil {
			log.Panicln("Can't hash password: ", err)
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		err = qtx.UpdateUserPasswordHash(r.Context(), db.UpdateUserPasswordHashParams{
			ID:           passwordResetToken.UserId,
			PasswordHash: passwordHash,
		})
		if err != nil {
			log.Panicln("Can't update password: " + err.Error())
		}

		err = qtx.DeletePasswordResetToken(r.Context(), tokenHash)
		if err != nil {
			log.Panicln("Can't delete password reset token: " + err.Error())
		}

		// Invalidate all of the user's sessions
		err = app.Scs.Iterate(r.Context(), func(ctx context.Context) error {
			userID := app.Scs.GetBytes(ctx, "userId")
			if len(userID) == 16 && [16]byte(userID) == passwordResetToken.UserId.Bytes {
				return app.Scs.Destroy(ctx)
			}
			return nil
		})
		if err != nil {
			log.Panicln("Can't destroy sessions: ", err)
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
		href := r.Header.Get("Origin") + "/reset-password"

		if lng == "my" {
			templateData = struct {
				Subject string
				Body    string
			}{
				Subject: "သင့်လျှိဝှက်နံပါတ်အားပြန်လည်တပ်ဆင်ပြီးပါပြီ။",
				Body: `
					<p>Hello ` + passwordResetToken.Username + `<p/>
					<p>သင့်အကောင့်အားလျှိဝှက်နံပါတ်ပြန်လည်သတ်မှတ်ပြီးပါပြီ။ အကယ်၍အဲ့တာကသင်ဖြစ်လျှင်ဒီ  email ကိုလုံခြုံစွာလျစ်လျှူရှုနိုင်ပါသည်။<p/>
					<p>အကယ်၍သင်မဟုတ်လျှင် သင့်အကောင့်ကစိတ်မမချရဖြစ်နိုင်ပါသည်။ ချက်ချင်းသင့်အကောင့်ကိုလုံခြုံစေရန် <a href="` + href + "\">" + href + `</a>အားနှိပ်ပါ။‌<p/>
					<p>ပျာ်ရွှင်ပါစေ:))</p>`,
			}
		} else {
			templateData = struct {
				Subject string
				Body    string
			}{
				Subject: "Your password has been reset",
				Body: `
					<p>Hello ` + passwordResetToken.Username + `,<p/>
					<p>The password to your account has just been reset. If this was you, you can safely ignore this email.<p/>
					<p>If this wasn't you, your account may be compromised. Please secure your account by visiting <a href="` + href + "\">" + href + `</a> immediately.</p>
					<p>Cheers.</p>`,
			}
		}

		body, err := utils.ParseTemplate("template.html", templateData)
		if err != nil {
			log.Panicln("Cannot parse template: ", err.Error())
		}

		go func() {
			err := mailutils.SendMail(passwordResetToken.Email, body, templateData.Subject)
			if err != nil {
				log.Println("Can't send mail: " + err.Error())
			}
		}()

		utils.LogEvent(qtx, r, passwordResetToken.UserId, db.EventTypePASSWORDRESET, db.EventResultSUCCESS, "", nil)
		if err != nil {
			log.Panicln("Can't log event: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}

func checkPasswordResetToken(app *app.App, r *http.Request, tokenHash string) (db.GetPasswordResetTokenByHashRow, error) {
	key := "password_reset_ip:" + r.RemoteAddr
	err := app.Limiter.Consume(r.Context(), key, passwordResetIpLimitOpts)
	if err == ratelimiter.RateLimited {
		err2 := utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypePASSWORDRESETTOKENVERIFICATION, db.EventResultFAIL, "Rate limited", map[string]any{
			"tokenHash": tokenHash,
		})
		if err2 != nil {
			log.Panicln("Can't log event: " + err2.Error())
		}

		return db.GetPasswordResetTokenByHashRow{}, err
	}

	passwordResetToken, err := app.DB.GetPasswordResetTokenByHash(r.Context(), tokenHash)
	if err != nil || passwordResetToken.TokenHash != tokenHash {
		err2 := utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypePASSWORDRESETTOKENVERIFICATION, db.EventResultFAIL, "Password reset link invalid", map[string]any{
			"tokenHash": tokenHash,
		})
		if err2 != nil {
			log.Panicln("Can't log event: " + err2.Error())
		}
		return db.GetPasswordResetTokenByHashRow{}, err
	}

	expired := !time.Now().Before(passwordResetToken.CreatedAt.Time.Add(1 * time.Hour))
	if expired {
		err2 := utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypePASSWORDRESETTOKENVERIFICATION, db.EventResultFAIL, "Password reset link expired", map[string]any{
			"tokenHash": tokenHash,
		})
		if err2 != nil {
			log.Panicln("Can't log event: " + err2.Error())
		}
		return db.GetPasswordResetTokenByHashRow{}, err
	}

	return passwordResetToken, nil
}
