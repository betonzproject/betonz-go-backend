package producttype

import (
	"net/http"

	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
	"github.com/go-chi/chi/v5"
)

type ProductInfo struct {
	Name         string `json:"name"`
	UriComponent string `json:"uriComponent"`
	HasGameList  bool   `json:"hasGameList"`
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	productTypeParam := chi.URLParam(r, "productType")
	productType := product.UriComponentToProductType(productTypeParam)

	f := func(p product.Product) ProductInfo {
		return ProductInfo{
			Name:         p.String(),
			UriComponent: p.UriComponent(),
			HasGameList:  product.HasGameList(productType, p),
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
