package routes

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

type RegisterForm struct {
	Username  string `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Email     string `form:"email" validate:"required,email" key:"user.email"`
	Password  string `form:"password" validate:"required,min=8,max=512"`
	InvitedBy string `form:"invitedBy"`
}

var registerIpLimitOpts = ratelimiter.LimiterOptions{
	Tokens: 20,
	Window: time.Duration(24 * time.Hour),
}

func PostRegister(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerForm RegisterForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &registerForm) != nil {
			return
		}

		key := "register_ip:" + r.RemoteAddr
		err := app.Limiter.Consume(r.Context(), key, registerIpLimitOpts)
		if err == ratelimiter.RateLimited {
			err := utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypeREGISTER, db.EventResultFAIL, "Rate limited", map[string]any{
				"username": registerForm.Username,
				"email":    registerForm.Email,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			http.Error(w, "tooManyRequests.message", http.StatusTooManyRequests)
			return
		}

		_, err = app.DB.GetExtendedUserByUsername(r.Context(), db.GetExtendedUserByUsernameParams{
			Username: registerForm.Username,
		})
		if err == nil {
			err = utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypeREGISTER, db.EventResultFAIL, "Username already taken", map[string]any{
				"username": registerForm.Username,
				"email":    registerForm.Email,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			http.Error(w, "user.username.alreadyTaken.message", http.StatusForbidden)
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		passwordHash, _ := utils.Argon2IDHash(registerForm.Password)
		SendEmailVerification(qtx, r, db.User{}, &db.RegisterInfo{
			Username:     registerForm.Username,
			Email:        registerForm.Email,
			PasswordHash: passwordHash,
			InvitedBy:    registerForm.InvitedBy,
		})

		err = utils.LogEvent(qtx, r, pgtype.UUID{}, db.EventTypeREGISTER, db.EventResultSUCCESS, "", map[string]any{
			"username": registerForm.Username,
			"email":    registerForm.Email,
		})
		if err != nil {
			log.Panicln("Can't log event: ", err)
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
	}
}

func SendEmailVerification(q *db.Queries, r *http.Request, user db.User, registerInfo *db.RegisterInfo) {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	token := base64.RawURLEncoding.EncodeToString(randomBytes)

	hash := sha256.New()
	hash.Write([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

	err := q.UpsertVerificationToken(r.Context(), db.UpsertVerificationTokenParams{
		TokenHash:    tokenHash,
		UserId:       user.ID,
		RegisterInfo: registerInfo,
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

	var username, email string
	if registerInfo != nil {
		username = registerInfo.Username
		email = registerInfo.Email
	} else {
		username = user.Username
		email = user.Email
	}

	if lng == "my" {
		templateData = struct {
			Subject string
			Body    string
		}{
			Subject: "အီးမေးအတည်ပြု",
			Body: `
				<p>Hello ` + username + `,</p>
				<p>Beton မှ လှိက်လှဲစွာကြိုဆိုပါတယ်</p>
				<p>အကောင့်ဖွင့်ခြင်းအား ပီးမြောက်စေပီး game များစတင်ကစားနိုင်ရန်အတွက် 
				email အတည်ပြုရန်သာကျန်ပါတော့သည်။ email အတည်ပြုခြင်းလင့် သည် 24 
				နာရီကြာပီးချိန်တွင် သက်တမ်းကုန်ပါမည်။</p>
				<center style="margin-top: 10px;"><button style="color:black;background:#f3b83d;padding:.5rem .8rem;border-radius:999px;border:none"><a style="color:black;text-decoration:none" href="` + href + "\">Verify Email</a></button></center>",
		}
	} else {
		templateData = struct {
			Subject string
			Body    string
		}{
			Subject: "Verify Email",
			Body: `
				<p>Hello ` + username + `,<p />
				<p>Welcome to BetOn! We're thrilled to have you join our community.</p>
				<p>To complete your registration and gain access to our games, exclusive offers, and 24-hour customer service, we just need you to verify your email address. 
				Click the button below to verify your email. The link will expire in 1 hour.<p/>
				<center style="margin-top: 10px;"><button style="color:black;background:#f3b83d;padding:.5rem .8rem;border-radius:999px;border:none"><a style="color:black;text-decoration:none" href="` + href + "\">Verify Email</a></button></center>",
		}
	}

	body, err := utils.ParseTemplate("template.html", templateData)
	if err != nil {
		log.Panicln("Can't parse template : ", err.Error())
	}

	go func() {
		err := mailutils.SendMail(email, body, templateData.Subject)
		log.Println("Email sent!")
		if err != nil {
			log.Println("Can't send mail: " + err.Error())
		}
	}()
}
