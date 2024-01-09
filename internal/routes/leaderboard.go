package routes

import (
	"net/http"
	"strconv"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ranking struct {
	Id           int32          `json:"id"`
	Name         string         `json:"name"`
	ProfileImage pgtype.Text    `json:"profileImage"`
	Amount       pgtype.Numeric `json:"amount"`
}

func GetLeaderboard(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productTypeParam := chi.URLParam(r, "productType")
		productType, _ := strconv.Atoi(productTypeParam)

		payout, _ := app.DB.GetTopPayout(r.Context(), int32(productType))

		rankings := make([]ranking, 0)
		for _, p := range payout {
			name := p.DisplayName.String
			if name == "" {
				name = p.Username
			}

			nameToShow := "*****" + name[len(name)-3:]
			if len(name) > 4 {
				nameToShow = "*****" + name[len(name)-3:]
			} else {
				nameToShow = "*****" + name[:len(name)-1]
			}

			ranking := ranking{
				Id:           p.ID,
				Name:         nameToShow,
				ProfileImage: p.ProfileImage,
				Amount:       p.Payout,
			}
			rankings = append(rankings, ranking)
		}

		jsonutils.Write(w, rankings, http.StatusOK)
	}
}
