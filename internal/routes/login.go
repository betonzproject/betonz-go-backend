package routes

import (
	"crypto/rand"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/mailutils"
	"github.com/doorman2137/betonz-go/internal/utils/ratelimiter"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetLogin(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.GetUser(app, w, r)
		if err == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		jsonutils.Write(w, struct{}{}, http.StatusOK)
	}
}

type LoginForm struct {
	Username string `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Password string `form:"password" validate:"required,min=8,max=512"`
	Pin      string `form:"pin"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Pin   bool   `json:"pin"`
}

var loginIpLimitOpts = ratelimiter.LimiterOptions{
	Tokens: 100,
	Window: time.Duration(24 * time.Hour),
}

var loginIpUsernameLimitOpts = ratelimiter.LimiterOptions{
	Tokens: 20,
	Window: time.Duration(24 * time.Hour),
}

func PostLogin(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginForm LoginForm
		if formutils.ParseDecodeValidateMultipart(app, w, r, &loginForm) != nil {
			return
		}

		adminMode := r.URL.Query().Get("role") == string(db.RoleADMIN)

		redirectParam := r.URL.Query().Get("redirect")
		redirectTo, err := url.QueryUnescape(redirectParam)
		if err != nil || redirectTo == "" {
			redirectTo = "/"
		} else {
			redirectTo = "/" + redirectTo[1:]
		}

		ipKey := "login_ip:" + r.RemoteAddr
		ipUsernameKey := "login_ip_username:" + r.RemoteAddr + ":" + loginForm.Username
		err = app.Limiter.Consume(r.Context(), ipKey, loginIpLimitOpts)
		err2 := app.Limiter.Consume(r.Context(), ipUsernameKey, loginIpUsernameLimitOpts)
		if err == ratelimiter.RateLimited || err2 == ratelimiter.RateLimited {
			err := utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypeLOGIN, db.EventResultFAIL, "Rate limited", map[string]any{
				"redirectTo": redirectTo,
				"adminMode":  adminMode,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			http.Error(w, "login.tooManyAttempts.message", http.StatusTooManyRequests)
			return
		}

		var user db.User
		if adminMode {
			user, err = app.DB.GetExtendedUserByUsername(r.Context(), db.GetExtendedUserByUsernameParams{
				Username: loginForm.Username,
				Roles: []db.Role{
					db.RoleADMIN,
					db.RoleSUPERADMIN,
				},
			})
		} else {
			user, err = app.DB.GetExtendedUserByUsername(r.Context(), db.GetExtendedUserByUsernameParams{
				Username: loginForm.Username,
				Roles: []db.Role{
					db.RolePLAYER,
				},
			})
		}
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				log.Panicln("Can't get user: " + err.Error())
			}

			// Dummy hash to prevent timing attack
			utils.Argon2IDVerify(loginForm.Password, "$argon2id$v=19$m=65536,t=3,p=4$YGmXRJpAsWMPAU8eMrFFIw$P5NwS7fKyuGaU+siAOFeNBmfbucV3Rrj7rEUMSB4vc8")

			err := utils.LogEvent(app.DB, r, user.ID, db.EventTypeLOGIN, db.EventResultFAIL, "User does not exist", map[string]any{
				"username":   loginForm.Username,
				"redirectTo": redirectTo,
				"adminMode":  adminMode,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			http.Error(w, "login.usernameOrPasswordIncorrect.message", http.StatusUnauthorized)
			return
		}

		passwordMatches, _ := utils.Argon2IDVerify(loginForm.Password, user.PasswordHash)
		if !passwordMatches {
			err := utils.LogEvent(app.DB, r, user.ID, db.EventTypeLOGIN, db.EventResultFAIL, "Password incorrect", map[string]any{
				"username":   loginForm.Username,
				"redirectTo": redirectTo,
				"adminMode":  adminMode,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			http.Error(w, "login.usernameOrPasswordIncorrect.message", http.StatusUnauthorized)
			return
		}

		if adminMode && loginForm.Pin == "" {
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

			code := utils.GeneratePIN(8)

			err = app.DB.UpsertVerificationPin(r.Context(), db.UpsertVerificationPinParams{
				Pin:    code,
				UserId: user.ID,
			})
			if err != nil {
				log.Panicln("Cannot create verification pin: ", err.Error())
			}

			if lng == "my" {
				templateData = struct {
					Subject string
					Body    string
				}{
					Subject: "Verification Code",
					Body: `<h3 style="color:white;">Account ဝင်ရင် အောက်ကကုဒ်ကို ဖြည့်ပါ</h3>
					<center>
						<h1 style="background: #f3b83d; border-radius: 1rem; color: black; padding: 1rem; letter-spacing: 1rem;">` + code + `</h1>
					</center>`,
				}
			} else {
				templateData = struct {
					Subject string
					Body    string
				}{
					Subject: "Verification Code",
					Body: `<h3 style="color:white;">Use the following code to access your account.</h3>
					<center>
						<h1 style="background: #f3b83d; border-radius: 1rem; color: black; padding: 1rem; letter-spacing: 1rem;">` + code + `</h1>
					</center>`,
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

			return
		} else if adminMode {
			verificationPIN, err := app.DB.GetVerificationPinByUserId(r.Context(), user.ID)
			if err != nil {
				return
			}

			expired := !time.Now().Before(verificationPIN.CreatedAt.Time.Add(10 * time.Minute))

			if verificationPIN.Pin != loginForm.Pin || expired {
				jsonutils.Write(w, LoginResponse{Pin: true}, http.StatusOK)
				return
			}

			err = app.DB.DeleteVerificationPin(r.Context(), verificationPIN.Pin)
			if err != nil {
				log.Panicln("Error deleting pin: ", err.Error())
			}
		}

		app.Scs.RenewToken(r.Context())
		app.Scs.Put(r.Context(), "userId", user.ID.Bytes[:])
		randomBytes := make([]byte, 16)
		rand.Read(randomBytes)
		app.Scs.Put(r.Context(), "sessionId", randomBytes)

		app.Limiter.Reset(r.Context(), ipUsernameKey)

		err = utils.LogEvent(app.DB, r, user.ID, db.EventTypeLOGIN, db.EventResultSUCCESS, "", map[string]any{
			"redirectTo": redirectTo,
			"adminMode":  adminMode,
		})
		if err != nil {
			log.Panicln("Can't log event: " + err.Error())
		}

		http.Redirect(w, r, redirectTo, http.StatusFound)
	}
}
