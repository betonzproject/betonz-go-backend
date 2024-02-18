package verifyemail

import (
	"crypto/sha256"
	"encoding/base64"
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
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTokenResponse struct {
	IsTokenValid bool `json:"isTokenValid"`
}

func GetVerifyEmailToken(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")

		hash := sha256.New()
		hash.Write([]byte(token))
		tokenHash := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

		emailVerificationToken, err := app.DB.GetVerificationTokenByHash(r.Context(), tokenHash)
		if err != nil || emailVerificationToken.TokenHash != tokenHash {
			err := utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypeEMAILVERIFICATION, db.EventResultFAIL, "Email verification link invalid", map[string]any{
				"tokenHash": tokenHash,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			jsonutils.Write(w, VerifyEmailTokenResponse{IsTokenValid: false}, http.StatusBadRequest)
			return
		}

		expired := !time.Now().Before(emailVerificationToken.CreatedAt.Time.Add(1 * time.Hour))
		if expired {
			err := utils.LogEvent(app.DB, r, pgtype.UUID{}, db.EventTypeEMAILVERIFICATION, db.EventResultFAIL, "Email verification link expired", map[string]any{
				"tokenHash": tokenHash,
			})
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}

			jsonutils.Write(w, VerifyEmailTokenResponse{IsTokenValid: false}, http.StatusBadRequest)
			return
		}

		etgUsername, err := createPlayer()
		if err != nil {
			log.Panicln("Can't create player: ", err)
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		var userId pgtype.UUID
		registerInfo := emailVerificationToken.RegisterInfo
		if registerInfo != nil {
			user, err := qtx.CreateUser(r.Context(), db.CreateUserParams{
				Username:     registerInfo.Username,
				Email:        registerInfo.Email,
				PasswordHash: registerInfo.PasswordHash,
				EtgUsername:  etgUsername,
			})
			if err != nil {
				log.Panicln("Can't create new user: ", err)
			}

			userId = user.ID
		} else {
			err := qtx.MarkUserEmailAsVerified(r.Context(), emailVerificationToken.UserId)
			if err != nil {
				log.Panicln("Can't mark email as verified: " + err.Error())
			}

			userId = emailVerificationToken.UserId
		}

		err = qtx.DeleteVerificationTokenByHash(r.Context(), tokenHash)
		if err != nil {
			log.Panicln("Can't delete email verification token: ", err)
		}

		err = utils.LogEvent(qtx, r, userId, db.EventTypeEMAILVERIFICATION, db.EventResultSUCCESS, "", map[string]any{
			"tokenHash": tokenHash,
		})
		if err != nil {
			log.Panicln("Can't create event: " + err.Error())
		}

		tx.Commit(r.Context())

		jsonutils.Write(w, VerifyEmailTokenResponse{IsTokenValid: true}, http.StatusOK)
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
			log.Panicf("Can't create player: %s\nEndpoint: %s\nPayload: %+v\n", err, endpoint, payload)
		}

		if createPlayerResponse.Err != etg.Success {
			return "", fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %+v", createPlayerResponse.Err, createPlayerResponse.Desc, endpoint, payload)
		}
	}

	return string(etgUsername), nil
}
