package routes

import (
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
)

func PostLogout(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Scs.RenewToken(r.Context())
		app.Scs.Destroy(r.Context())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
