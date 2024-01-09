package routes

import (
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
)

type LoginForm struct {
	Username string `formam:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Password string `formam:"password" validate:"required,min=8,max=512"`
}

func PostLogin(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginForm LoginForm
		if formutils.ParseDecodeValidate(app, w, r, &loginForm) != nil {
			return
		}

		adminMode := r.URL.Query().Get("role") == string(db.RoleADMIN)

		var user db.User
		var err error
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
			// Dummy hash to prevent timing attack
			utils.Argon2IDVerify(loginForm.Password, "$argon2id$v=19$m=65536,t=3,p=4$YGmXRJpAsWMPAU8eMrFFIw$P5NwS7fKyuGaU+siAOFeNBmfbucV3Rrj7rEUMSB4vc8")
			http.Error(w, "login.usernameOrPasswordIncorrect.message", http.StatusUnauthorized)
			return
		}

		passwordMatches, _ := utils.Argon2IDVerify(loginForm.Password, user.PasswordHash)
		if !passwordMatches {
			http.Error(w, "login.usernameOrPasswordIncorrect.message", http.StatusUnauthorized)
			return
		}

		app.Scs.RenewToken(r.Context())
		app.Scs.Put(r.Context(), "userId", user.ID.Bytes[:])

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
