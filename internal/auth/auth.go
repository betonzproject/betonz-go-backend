package auth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// Checks that the user is authenticated and returns the user. Otherwise, redirect the user to
// the login page and an error is returned
func Authenticate(app *app.App, w http.ResponseWriter, r *http.Request) (db.User, error) {
	userId := app.Scs.GetBytes(r.Context(), "userId")

	redirectTo := url.QueryEscape(r.URL.Path)

	// Strip /admin route prefix for admin routes
	if strings.HasPrefix(redirectTo, "/admin") {
		redirectTo = redirectTo[6:]
	}

	if len(userId) < 16 {
		http.Redirect(w, r, "/login?redirect="+url.QueryEscape(redirectTo), http.StatusFound)
		return db.User{}, errors.New("Unauthenticated")
	}

	user, err := app.DB.GetExtendedUserById(r.Context(), pgtype.UUID{Bytes: [16]byte(userId), Valid: true})
	if err != nil {
		http.Redirect(w, r, "/login?redirect="+url.QueryEscape(redirectTo), http.StatusFound)
		return user, err
	}

	return user, nil
}
