package producttype

import (
	"log"
	"net/http"
	"slices"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
	"github.com/go-chi/chi/v5"
)

type ProductInfo struct {
	Name               string `json:"name"`
	UriComponent       string `json:"uriComponent"`
	HasGameList        bool   `json:"hasGameList"`
	IsUnderMaintenance bool   `json:"isUnderMaintenance"`
}

func GetProducts(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productTypeParam := chi.URLParam(r, "productType")
		productType := product.UriComponentToProductType(productTypeParam)

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		f := func(p product.Product) ProductInfo {
			isUnderMaintenance := false
			if slices.Contains(productsUnderMaintenance, int32(p)) {
				isUnderMaintenance = true
			}

			return ProductInfo{
				Name:               p.String(),
				UriComponent:       p.UriComponent(),
				HasGameList:        product.HasGameList(productType, p),
				IsUnderMaintenance: isUnderMaintenance,
			}
		}

		switch productType {
		case product.LiveCasino:
			jsonutils.Write(w, sliceutils.Map(product.LiveCasinoProducts, f), http.StatusOK)
		case product.Slots:
			jsonutils.Write(w, sliceutils.Map(product.SlotsProducts, f), http.StatusOK)
		case product.Sports:
			jsonutils.Write(w, sliceutils.Map(product.SportsProducts, f), http.StatusOK)
		case product.CardAndBoard:
			jsonutils.Write(w, sliceutils.Map(product.CardAndBoardProducts, f), http.StatusOK)
		case product.Fishing:
			jsonutils.Write(w, sliceutils.Map(product.FishingProducts, f), http.StatusOK)
		case product.Cockfighting:
			jsonutils.Write(w, sliceutils.Map(product.CockfightingProducts, f), http.StatusOK)
		default:
			http.Error(w, "404 page not found", http.StatusNotFound)
		}
	}
}
