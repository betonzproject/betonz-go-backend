package transfer

import (
	"net/http"

	"github.com/doorman2137/betonz-go/internal/product"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
)

type Products struct {
	Products     map[product.Product]string `json:"products"`
	WalletGroups [][]product.Product        `json:"walletGroups"`
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	productNames := make(map[product.Product]string)

	for _, p := range product.AllProducts {
		productNames[p] = p.String()
	}

	jsonutils.Write(w, Products{
		Products:     productNames,
		WalletGroups: [][]product.Product{product.WalletGroup1, product.WalletGroup2},
	}, http.StatusOK)
}
