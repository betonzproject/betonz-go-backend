package resetpassword

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/mailutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type ResetPasswordRequestForm struct {
	Username string `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Email    string `form:"email" validate:"required,email" key:"user.email"`
}

func PostResetPassword(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resetPasswordRequestForm ResetPasswordRequestForm
		if formutils.ParseDecodeValidate(app, w, r, &resetPasswordRequestForm) != nil {
			return
		}

		tx, err := app.Pool.Begin(r.Context())
		if err != nil {
			log.Panicln("Can't start transaction: " + err.Error())
		}
		defer tx.Rollback(r.Context())
		qtx := app.DB.WithTx(tx)

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

		exisitingToken, err := qtx.GetPasswordResetTokenByUserId(r.Context(), user.ID)
		if err == nil {
			err = qtx.DeletePasswordResetToken(r.Context(), exisitingToken.TokenHash)
			if err != nil {
				log.Panicln("Can't delete password reset token: " + err.Error())
			}
		}

		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		token := base64.RawURLEncoding.EncodeToString(randomBytes)

		hash := sha256.New()
		hash.Write([]byte(token))
		tokenHash := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

		err = qtx.CreatePasswordResetToken(r.Context(), db.CreatePasswordResetTokenParams{
			TokenHash: tokenHash,
			UserId:    user.ID,
		})
		if err != nil {
			log.Panicln("Cannot create password reset token: ", err.Error())
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
					<center><button style="color:white;background:#f3b83d;padding:.5rem .8rem;border-radius:999px;border:none"><a style="color:black;text-decoration:none" href="` + href + "\">ပစ်စဝက်စ် ပြန်ပေးမည်</a></button></center>",
			}
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
					<center><button style="color:white;background:#f3b83d;padding:.5rem .8rem;border-radius:999px;border:none"><a style="color:black;text-decoration:none" href="` + href + "\">Reset Password</a></button></center>",
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
