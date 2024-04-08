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
	"github.com/doorman2137/betonz-go/internal/utils/ratelimiter"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
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
}

type LoginResponse struct {
	Email string `json:"email"`
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
		if formutils.ParseDecodeValidate(app, w, r, &loginForm) != nil {
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

		// If the user is a player and email is not verified, ask for email verification
		if !adminMode && !user.IsEmailVerified {
			log.Println(user)
			tx, qtx := transactionutils.Begin(app, r.Context())
			defer tx.Rollback(r.Context())

			SendEmailVerification(qtx, r, user, nil)

			err := utils.LogEvent(qtx, r, user.ID, db.EventTypeLOGIN, db.EventResultFAIL, "Email is not verified", map[string]any{
				"redirectTo": redirectTo,
				"adminMode":  adminMode,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			tx.Commit(r.Context())

			jsonutils.Write(w, LoginResponse{user.Email}, http.StatusOK)
			return
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
