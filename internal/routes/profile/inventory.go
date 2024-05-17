package profile

import (
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
)

type InventoryResponse struct {
	Inventory []db.GetInventoryByUserIdRow `json:"inventory"`
}

func GetInventory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		userInventory, err := app.DB.GetInventoryByUserId(r.Context(), user.ID)
		if err != nil {
			return
		}

		jsonutils.Write(w, InventoryResponse{
			Inventory: userInventory,
		}, http.StatusOK)
	}
}
