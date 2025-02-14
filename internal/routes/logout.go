package routes

import (
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/app"
)

func PostLogout(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Scs.Destroy(r.Context())
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
