package resetpassword

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/formutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/mailutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/ratelimiter"
	"github.com/BetOnz-Company/betonz-go/internal/utils/transactionutils"

	"github.com/jackc/pgx/v5/pgtype"
)

type ResetPasswordRequestForm struct {
	Username string `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Email    string `form:"email" validate:"required,email" key:"user.email"`
}

var passwordResetIpLimitOpts = ratelimiter.LimiterOptions{
	Tokens: 20,
	Window: time.Duration(24 * time.Hour),
}

var passwordResetUsernameLimitOpts = ratelimiter.LimiterOptions{
	Tokens: 10,
	Window: time.Duration(24 * time.Hour),
}

func PostPasswordReset(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resetPasswordRequestForm ResetPasswordRequestForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &resetPasswordRequestForm) != nil {
			return
		}

		ipKey := "password_reset_ip:" + r.RemoteAddr
		ipUsernameKey := "password_reset_username:" + resetPasswordRequestForm.Username
		err := app.Limiter.Consume(r.Context(), ipKey, passwordResetIpLimitOpts)
		err2 := app.Limiter.Consume(r.Context(), ipUsernameKey, passwordResetUsernameLimitOpts)
		if err == ratelimiter.RateLimited || err2 == ratelimiter.RateLimited {
			err := utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypePASSWORDRESETREQUEST, db.EventResultFAIL, "Rate limited", map[string]any{
				"username": resetPasswordRequestForm.Username,
				"email":    resetPasswordRequestForm.Email,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			http.Error(w, "resetPassword.tooManyPasswordResetRequests.message", http.StatusTooManyRequests)
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		user, err := qtx.GetExtendedUserByUsername(r.Context(), db.GetExtendedUserByUsernameParams{
			Username: resetPasswordRequestForm.Username,
			Roles:    []db.Role{db.RolePLAYER},
		})

		if err != nil || resetPasswordRequestForm.Email != user.Email {
			err = utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypePASSWORDRESETREQUEST, db.EventResultFAIL, "Username or email does not match", map[string]any{
				"username": resetPasswordRequestForm.Username,
				"email":    resetPasswordRequestForm.Email,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		token := base64.RawURLEncoding.EncodeToString(randomBytes)

		hash := sha256.New()
		hash.Write([]byte(token))
		tokenHash := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

		err = qtx.UpsertPasswordResetToken(r.Context(), db.UpsertPasswordResetTokenParams{
			TokenHash: tokenHash,
			UserId:    user.ID,
		})
		if err != nil {
			log.Panicln("Cannot upsert password reset token: ", err.Error())
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
		href := r.Header.Get("Origin") + "/reset-password/" + token

		if lng == "my" {
			templateData = struct {
				Subject string
				Body    string
			}{
				Subject: "လျှိဝှက်နံပါတ်ပြန်လည်သတ်မှတ်",
				Body: `
					<p>Hello ` + user.Username + `,<p/>
					<p>သင််သည်လျှိဝှက်နံပါတ်ပြန်လည်သတ်မှတ်ရန်ခွင့်တောင်းထားပါသည်။</p>
					<p>လျှိဝှက်နံပါတ်ကိုပြန်လည်သတ်မှတ်ရန်အောက်မှာချပေးထားတဲ့linkကိုနှိပ်ပေးပါ။</p>
					<p>linkသည်၁နာရီအတွင်းတွင်သက်တမ်းကုန်မည်။ သင်သည်လျှိဝှက်နံပါတ်ပြန်လည်သတ်မှတ်ဖိုမတောင်းဆိုထားလျှင်, ဒီemailကိုလျစ်လျူရှုပေးပါ။ <p/>
					<center style="margin-top: 10px;"><a style="color:black;background:#f3b83d;padding:.5rem .8rem;border-radius:999px;border:none;text-decoration:none;display:inline-block;" href="` + href + `">Verify Email</a></center>`}
		} else {
			templateData = struct {
				Subject string
				Body    string
			}{
				Subject: "Password Reset",
				Body: `
					<p>Hello ` + user.Username + `,<p/>
					<p>You have requested to reset your password. Please click the link below to reset your password. The link will expire in 1 hour.<p/>
					<p>If you didn't request to reset your password, please ignore this email.<p/>
					<center style="margin-top: 10px;"><a style="color:black;background:#f3b83d;padding:.5rem .8rem;border-radius:999px;border:none;text-decoration:none;display:inline-block;" href="` + href + `">Verify Email</a></center>`,
			}
		}

		body, err := utils.ParseTemplate("template.html", templateData)
		if err != nil {
			log.Panicln("Can't parse template: ", err.Error())
		}

		go func() {
			err := mailutils.SendMail(resetPasswordRequestForm.Email, body, templateData.Subject)
			if err != nil {
				log.Println("Can't send mail: " + err.Error())
			}
		}()

		err = utils.LogEvent(qtx, r, user.ID, db.EventTypePASSWORDRESETREQUEST, db.EventResultSUCCESS, "", map[string]any{
			"username": resetPasswordRequestForm.Username,
			"email":    resetPasswordRequestForm.Email,
		})
		if err != nil {
			log.Panicln("Can't log event: " + err.Error())
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}
