package auth

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// Checks that the user is authenticated and returns the user. Otherwise, redirect the user to
// the login page and an error is returned
func Authenticate(app *app.App, w http.ResponseWriter, r *http.Request, redirect string) (db.GetUserByIdRow, error) {
	userId := app.Scs.GetBytes(r.Context(), "userId")
	if len(userId) < 16 {
		if redirect != "" {
			http.Redirect(w, r, "/login?redirect="+url.QueryEscape(redirect), http.StatusFound)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return db.GetUserByIdRow{}, errors.New("Unauthenticated")
	}

	user, err := app.DB.GetUserById(r.Context(), pgtype.UUID{Bytes: [16]byte(userId), Valid: true})
	if err != nil {
		if redirect != "" {
			http.Redirect(w, r, "/login?redirect="+url.QueryEscape(redirect), http.StatusFound)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return user, err
	}

	return user, nil
}
