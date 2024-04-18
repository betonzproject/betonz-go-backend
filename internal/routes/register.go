package routes

import (
	crand "crypto/rand"
	"fmt"
	"log"
	"math/rand/v2"

	"net/http"
	"os"
	"time"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/etg"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/ratelimiter"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegisterForm struct {
	Username string `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Email    string `form:"email" validate:"required,email" key:"user.email"`
	Password string `form:"password" validate:"required,min=8,max=512"`
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
		etgUsername, err := createPlayer()
		if err != nil {
			log.Panicln("Can't create player: ", err)
		}

		user, err := qtx.CreateUser(r.Context(), db.CreateUserParams{
			Username:     registerForm.Username,
			Email:        registerForm.Email,
			PasswordHash: passwordHash,
			EtgUsername:  etgUsername,
		})
		if err != nil {
			log.Panicln("Can't create new user: ", err)
		}

		err = utils.LogEvent(qtx, r, pgtype.UUID{}, db.EventTypeREGISTER, db.EventResultSUCCESS, "", map[string]any{
			"username": registerForm.Username,
			"email":    registerForm.Email,
		})
		if err != nil {
			log.Panicln("Can't log event: ", err)
		}

		app.Scs.RenewToken(r.Context())
		app.Scs.Put(r.Context(), "userId", user.ID.Bytes[:])
		randomBytes := make([]byte, 16)
		crand.Read(randomBytes)
		app.Scs.Put(r.Context(), "sessionId", randomBytes)

		tx.Commit(r.Context())

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

type CreatePlayerRequest struct {
	Op   string `json:"op"`
	Mem  string `json:"mem"`
	Pass string `json:"pass"`
}

type CreatePlayerResponse struct {
	Err  int    `json:"err"`
	Desc string `json:"desc"`
}

func createPlayer() (string, error) {
	// Generate random ETG username
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	etgUsername := make([]byte, 12)
	for i := range etgUsername {
		etgUsername[i] = charset[rand.IntN(len(charset))]
	}

	if os.Getenv("DISABLE_ETG_CREATE_PLAYER") == "" {
		endpoint := os.Getenv("ETG_API_ENDPOINT") + "/createplayer"

		payload := CreatePlayerRequest{
			Op:   os.Getenv("ETG_OPERATOR_CODE"),
			Mem:  string(etgUsername),
			Pass: "00000000",
		}
		var createPlayerResponse CreatePlayerResponse
		err := etg.Post("/createplayer", payload, &createPlayerResponse)
		if err != nil {
			log.Panicln("Can't create player: " + err.Error())
		}

		if createPlayerResponse.Err != etg.Success {
			return "", fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %+v", createPlayerResponse.Err, createPlayerResponse.Desc, endpoint, payload)
		}
	}

	return string(etgUsername), nil
}
