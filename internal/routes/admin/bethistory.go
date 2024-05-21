package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/product"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/sliceutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/timeutils"

	"github.com/jackc/pgx/v5/pgtype"
)

type Bet struct {
	ID               int32              `json:"id"`
	Username         string             `json:"username"`
	Role             db.Role            `json:"role"`
	RefId            string             `json:"refId"`
	EtgUsername      string             `json:"etgUsername"`
	ProviderUsername string             `json:"providerUsername"`
	Product          string             `json:"product"`
	ProductType      string             `json:"productType"`
	GameId           string             `json:"gameId"`
	Details          string             `json:"details"`
	Turnover         pgtype.Numeric     `json:"turnover"`
	Bet              pgtype.Numeric     `json:"bet"`
	Payout           pgtype.Numeric     `json:"payout"`
	Status           int32              `json:"status"`
	StartTime        pgtype.Timestamptz `json:"startTime"`
	MatchTime        pgtype.Timestamptz `json:"matchTime"`
	EndTime          pgtype.Timestamptz `json:"endTime"`
	SettleTime       pgtype.Timestamptz `json:"settleTime"`
	ProgShare        pgtype.Numeric     `json:"progShare"`
	ProgWIn          pgtype.Numeric     `json:"progWin"`
	Comission        pgtype.Numeric     `json:"comission"`
	WinLoss          pgtype.Numeric     `json:"winLoss"`
}

type BetHistoryResponse struct {
	ProductTypes map[product.ProductType]string `json:"productTypes"`
	Products     map[product.Product]string     `json:"products"`
	Bets         []Bet                          `json:"bets"`
}

func GetBetHistory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewBetHistory) != nil {
			return
		}

		searchParam := r.URL.Query().Get("search")
		fromParam := r.URL.Query().Get("from")
		toParam := r.URL.Query().Get("to")
		productCodeParam := r.URL.Query().Get("productCode")
		productTypeParam := r.URL.Query().Get("productType")

		var productCode int32
		if productCodeParam != "" {
			parsedProductCode, _ := strconv.Atoi(productCodeParam)
			productCode = int32(parsedProductCode)
		}

		var productType int32
		if productTypeParam != "" {
			parsedProductTypeproductType, _ := strconv.Atoi(productTypeParam)
			productType = int32(parsedProductTypeproductType)
		}

		from, err := timeutils.ParseDatetime(fromParam)
		if err != nil {
			from = timeutils.StartOfToday()
		}
		to, err := timeutils.ParseDatetime(toParam)
		if err != nil {
			to = timeutils.EndOfToday()
		}

		bets, err := app.DB.GetBets(r.Context(), db.GetBetsParams{
			Search:      pgtype.Text{String: searchParam, Valid: searchParam != ""},
			FromDate:    pgtype.Timestamptz{Time: from, Valid: true},
			ToDate:      pgtype.Timestamptz{Time: to, Valid: true},
			ProductCode: productCode,
			ProductType: productType,
		})
		if err != nil {
			log.Panicln("Can't fetch bets: ", err.Error())
		}

		productNames := make(map[product.Product]string)
		for _, p := range product.AllProducts {
			productNames[p] = p.String()
		}

		productTypes := make(map[product.ProductType]string)
		for _, p := range product.AllProductTypes {
			productTypes[p] = p.String()
		}

		jsonutils.Write(w,
			BetHistoryResponse{
				Products:     productNames,
				ProductTypes: productTypes,
				Bets: sliceutils.Map(bets, func(r db.GetBetsRow) Bet {
					betProduct := product.Product(int(r.ProductCode)).String()
					productType := product.ProductType(int(r.ProductType)).String()

					return Bet{
						ID:               r.ID,
						Username:         r.Username,
						Role:             r.Role,
						RefId:            r.RefId,
						EtgUsername:      r.EtgUsername,
						ProviderUsername: r.ProviderUsername,
						Product:          betProduct,
						ProductType:      productType,
						GameId:           r.GameId.String,
						Details:          r.Details,
						Turnover:         r.Turnover,
						Bet:              r.Bet,
						Payout:           r.Payout,
						Status:           r.Status,
						StartTime:        r.StartTime,
						MatchTime:        r.MatchTime,
						EndTime:          r.EndTime,
						SettleTime:       r.SettleTime,
						ProgShare:        r.ProgShare,
						ProgWIn:          r.ProgWin,
						Comission:        r.Commission,
						WinLoss:          r.WinLoss,
					}
				}),
			}, http.StatusOK)
	}
}
