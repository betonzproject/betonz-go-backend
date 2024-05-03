package producttype

import (
	"log"
	"net/http"
	"slices"
	"sort"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils/formutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetProductResponse struct {
	ProductName             string             `json:"productName"`
	Balance                 pgtype.Numeric     `json:"balance"`
	Games                   []product.GameInfo `json:"games"`
	ProductUnderMaintenance bool               `json:"productUnderMaintenance"`
}

func GetProduct(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productTypeParam := chi.URLParam(r, "productType")
		productType := product.UriComponentToProductType(productTypeParam)
		productParam := chi.URLParam(r, "product")
		p := product.UriComponentToProduct(productParam)

		if productType == 0 {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		if slices.Contains(productsUnderMaintenance, int32(p)) {
			jsonutils.Write(w, GetProductResponse{
				ProductName:             p.String(),
				ProductUnderMaintenance: true,
			}, http.StatusOK)
			return
		}

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

		var balance pgtype.Numeric
		user, err := auth.GetExtendedUser(app, w, r)
		if err == nil {
			balance, err = product.GetUserBalance(user.EtgUsername, p)
			if err != nil {
				log.Panicf("Can't get balance of %s (%d) for %s: %s\n", p, p, user.EtgUsername, err)
			}
		}

		games, err := product.GetGameList(app.Redis, r.Context(), productType, p)
		if err != nil {
			log.Panicln("Can't get games: " + err.Error())
		}

		if p == product.PragmaticPlay {
			sort.Slice(games, func(i, j int) bool {
				// Check if games[i] and games[j] are in PragmaticPlayTopGame
				iIndex := -1
				jIndex := -1
				for idx, game := range product.PragmaticPlayTopGame {
					if games[i].Name == game {
						iIndex = idx
					}
					if games[j].Name == game {
						jIndex = idx
					}
				}

				// If both are in PragmaticPlayTopGame, compare their indexes
				if iIndex != -1 && jIndex != -1 {
					return iIndex < jIndex
				}

				// If only games[i] is in PragmaticPlayTopGame, prioritize it
				if iIndex != -1 {
					return true
				}

				// If only games[j] is in PragmaticPlayTopGame, prioritize it
				if jIndex != -1 {
					return false
				}

				// If none of them are in PragmaticPlayTopGame, compare their names
				return games[i].Name < games[j].Name
			})
		}

		jsonutils.Write(w, GetProductResponse{
			ProductName:             p.String(),
			Games:                   games,
			Balance:                 balance,
			ProductUnderMaintenance: false,
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
		if formutils.ParseDecodeValidateMultipart(app, w, r, &gameForm) != nil {
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
