package routes

import (
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
)

type LoginForm struct {
	Username string `formam:"username" validate:"required"`
	Password string `formam:"password" validate:"required"`
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
			user, err = app.DB.GetExtendedAdminByUsername(r.Context(), loginForm.Username)
		} else {
			user, err = app.DB.GetExtendedPlayerByUsername(r.Context(), loginForm.Username)
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
