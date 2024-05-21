package auth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

// Returns the current authenticated user. Returns an error if the user is not authenticated.
func GetUser(app *app.App, w http.ResponseWriter, r *http.Request) (db.GetUserByIdRow, error) {
	userId := app.Scs.GetBytes(r.Context(), "userId")
	if len(userId) < 16 {
		return db.GetUserByIdRow{}, errors.New("Unauthenticated")
	}

	user, err := app.DB.GetUserById(r.Context(), pgtype.UUID{Bytes: [16]byte(userId), Valid: true})
	return user, err
}

// Returns the current authenticated user. Returns an error if the user is not authenticated.
//
// Same as `GetUser`, but returned user has more fields such as passwordHash, lastUsedBankId etc.
func GetExtendedUser(app *app.App, w http.ResponseWriter, r *http.Request) (db.User, error) {
	userId := app.Scs.GetBytes(r.Context(), "userId")
	if len(userId) < 16 {
		return db.User{}, errors.New("Unauthenticated")
	}

	user, err := app.DB.GetExtendedUserById(r.Context(), pgtype.UUID{Bytes: [16]byte(userId), Valid: true})
	return user, err
}

// Checks that the user is authenticated and returns the user. Otherwise, redirect the user to
// the login page and return an error
func Authenticate(app *app.App, w http.ResponseWriter, r *http.Request) (db.User, error) {
	user, err := GetExtendedUser(app, w, r)
	if err != nil {
		redirectTo := url.QueryEscape(r.URL.Path)

		// Strip /admin route prefix for admin routes
		if strings.HasPrefix(redirectTo, "/admin") {
			redirectTo = redirectTo[6:]
		}

		http.Redirect(w, r, "/login?redirect="+url.QueryEscape(redirectTo), http.StatusFound)
	}
	return user, err
}
