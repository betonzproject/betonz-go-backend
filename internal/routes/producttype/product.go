package producttype

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetProductResponse struct {
	ProductName string             `json:"productName"`
	Balance     float64            `json:"balance"`
	Games       []product.GameInfo `json:"games"`
}

func GetProduct(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productTypeParam := chi.URLParam(r, "productType")
		productType := product.UriComponentToProductType(productTypeParam)
		productParam := chi.URLParam(r, "product")
		p := product.UriComponentToProduct(productParam)

		if !product.HasGameList(productType, p) {
			user, err := auth.Authenticate(app, w, r)
			if err != nil {
				return
			}

			url, err := product.LaunchGameList(user.EtgUsername, productType, p)
			if err != nil {
				log.Panicln("Can't launch game list: " + err.Error())
			}

			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		var balance float64
		userId := app.Scs.GetBytes(r.Context(), "userId")
		if len(userId) == 16 {
			user, err := app.DB.GetExtendedUserById(r.Context(), pgtype.UUID{Bytes: [16]byte(userId), Valid: true})
			if err == nil {
				balance, err = product.GetUserBalance(user.EtgUsername, p)
				if err != nil {
					log.Panicf("Can't get balance of %s (%d) for %s: %s\n", p, p, user.EtgUsername, err)
				}
			}
		}

		games, err := product.GetGameList(app, r.Context(), productType, p)
		if err != nil {
			log.Panicln("Can't get games: " + err.Error())
		}

		jsonutils.Write(w, GetProductResponse{
			ProductName: p.String(),
			Games:       games,
			Balance:     balance,
		}, http.StatusOK)
	}
}

type GameForm struct {
	GameId string `form:"gameId" validate:"required"`
}

func PostProduct(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		var gameForm GameForm
		if formutils.ParseDecodeValidate(app, w, r, &gameForm) != nil {
			return
		}

		productTypeParam := chi.URLParam(r, "productType")
		productType := product.UriComponentToProductType(productTypeParam)
		productParam := chi.URLParam(r, "product")
		p := product.UriComponentToProduct(productParam)

		url, err := product.LaunchGame(user.EtgUsername, productType, p, gameForm.GameId)
		if err != nil {
			log.Panicln("Can't launch game: " + err.Error())
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
