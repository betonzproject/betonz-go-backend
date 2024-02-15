package routes

import (
	"crypto/rand"
	"net/http"
	"net/url"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
)

type LoginForm struct {
	Username string `form:"username" validate:"required,min=3,max=20,username" key:"user.username"`
	Password string `form:"password" validate:"required,min=8,max=512"`
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
			// Dummy hash to prevent timing attack
			utils.Argon2IDVerify(loginForm.Password, "$argon2id$v=19$m=65536,t=3,p=4$YGmXRJpAsWMPAU8eMrFFIw$P5NwS7fKyuGaU+siAOFeNBmfbucV3Rrj7rEUMSB4vc8")
			utils.LogEvent(app.DB, r, user.ID, db.EventTypeLOGIN, db.EventResultFAIL, "User does not exist", map[string]any{
				"username":   loginForm.Username,
				"redirectTo": redirectTo,
				"adminMode":  adminMode,
			})
			http.Error(w, "login.usernameOrPasswordIncorrect.message", http.StatusUnauthorized)
			return
		}

		passwordMatches, _ := utils.Argon2IDVerify(loginForm.Password, user.PasswordHash)
		if !passwordMatches {
			utils.LogEvent(app.DB, r, user.ID, db.EventTypeLOGIN, db.EventResultFAIL, "Password does not match", map[string]any{
				"username":   loginForm.Username,
				"redirectTo": redirectTo,
				"adminMode":  adminMode,
			})
			http.Error(w, "login.usernameOrPasswordIncorrect.message", http.StatusUnauthorized)
			return
		}

		app.Scs.RenewToken(r.Context())
		app.Scs.Put(r.Context(), "userId", user.ID.Bytes[:])
		randomBytes := make([]byte, 16)
		rand.Read(randomBytes)
		app.Scs.Put(r.Context(), "sessionId", randomBytes)

		utils.LogEvent(app.DB, r, user.ID, db.EventTypeLOGIN, db.EventResultSUCCESS, "", map[string]any{
			"redirectTo": redirectTo,
			"adminMode":  adminMode,
		})

		http.Redirect(w, r, redirectTo, http.StatusFound)
	}
}
