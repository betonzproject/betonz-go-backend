package transfer

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
)

type Products struct {
	Products                 map[product.Product]string `json:"products"`
	WalletGroups             [][]product.Product        `json:"walletGroups"`
	ProductsUnderMaintenance []string                   `json:"productsUnderMaintenance"`
}

func GetProducts(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productNames := make(map[product.Product]string)

		for _, p := range product.AllProducts {
			productNames[p] = p.String()
		}

		productsUnderMaintenance, err := app.DB.GetMaintenanceProductCodes(r.Context())
		if err != nil {
			log.Panicln("Error fetching maintained products: ", err.Error())
		}

		jsonutils.Write(w, Products{
			Products:     productNames,
			WalletGroups: [][]product.Product{product.WalletGroup1, product.WalletGroup2},
			ProductsUnderMaintenance: sliceutils.Map(productsUnderMaintenance, func(prodInt int32) string {
				return product.Product(prodInt).String()
			}),
		}, http.StatusOK)
	}
}
