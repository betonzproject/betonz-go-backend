package admin

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/etg"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
)

type CreateAdminForm struct {
	Username string `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Email    string `form:"email" validate:"required,email" key:"user.email"`
	Password string `form:"password" validate:"required,min=8,max=512"`
	Role     string `form:"role" validate:"oneof=ADMIN SUPERADMIN"`
}

type DeleteAdminForm struct {
	UserId string `form:"userId" validate:"uuid4"`
}

func GetAdmin(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageAdmins) != nil {
			return
		}

		players, err := app.DB.GetAdmins(r.Context(), []db.Role{db.RoleADMIN, db.RoleSUPERADMIN})
		if err != nil {
			log.Panicln("Can't get players: " + err.Error())
		}

		jsonutils.Write(w, players, http.StatusOK)
	}
}

func PostAdmin(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManageAdmins) != nil {
			return
		}

		tx, qtx := transactionutils.Begin(app, r.Context())
		defer tx.Rollback(r.Context())

		if r.URL.Query().Has("/create") {

			var createAdminForm CreateAdminForm
			if formutils.ParseDecodeValidateMultipart(app, w, r, &createAdminForm) != nil {
				return
			}

			_, err = qtx.GetExtendedUserByUsername(r.Context(), db.GetExtendedUserByUsernameParams{
				Username: createAdminForm.Username,
			})
			if err == nil {
				http.Error(w, "User already exist", http.StatusBadRequest)
				return
			}

			passwordHash, _ := utils.Argon2IDHash(createAdminForm.Password)
			etgUsername, err := createPlayer()
			if err != nil {
				log.Panicln("Error creating player for admin: ", err.Error())
			}

			_, err = qtx.CreateAdmin(r.Context(), db.CreateAdminParams{
				Username:     createAdminForm.Username,
				Email:        createAdminForm.Email,
				PasswordHash: passwordHash,
				Role:         db.Role(createAdminForm.Role),
				EtgUsername:  etgUsername,
			})
			if err != nil {
				log.Panicln("Error creating new admin: ", err.Error())
			}

		} else if r.URL.Query().Has("/delete") {
			var deleteAdminForm DeleteAdminForm
			if formutils.ParseDecodeValidateMultipart(app, w, r, &deleteAdminForm) != nil {
				return
			}

			userId, _ := utils.ParseUUID(deleteAdminForm.UserId)

			err = qtx.DeleteAdmin(r.Context(), userId)
			if err != nil {
				log.Panicln("Error deleting admin: ", err.Error())
			}
		}
		tx.Commit(r.Context())
		w.WriteHeader(http.StatusOK)

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
