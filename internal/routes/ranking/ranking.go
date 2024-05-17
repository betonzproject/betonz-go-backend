package ranking

import (
	"log"
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
)

func GetUserRanking(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		// Fetch User Ranking
		ranking, err := app.DB.GetUserRanking(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Error getting user ranking: ", err.Error())
			return
		}

		// Fetch Total Bet Amount
		totalBetAmount, err := app.DB.GetTotalBetAmount(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Error getting total bet amount: ", err.Error())
			return
		}

		weeklyTurnover, err := app.DB.GetWeeklyTurnover(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Error getting weekly turnover: ", err.Error())
			return
		}

		benefits, err := app.DB.GetUserBenefits(r.Context(), user.ID)
		if err != nil {
			log.Panicln("Error getting user benefits: ", err.Error())
			return
		}

		vipLevelStr := string(ranking.VipLevel.VipType)

		response := struct {
			VipLevel        string  `json:"vipLevel"`
			NewVipLevel     string  `json:"newVipLevel"`
			TotalBetAmount  int64   `json:"totalBetAmount"`
			WeeklyTurnover  int64   `json:"weeklyTurnover"`
			BirthdayBonus   int32   `json:"birthdayBonus"`
			BirthdayGift    string  `json:"birthdayGift"`
			MonthlyGift     string  `json:"monthlyGift"`
			AnniversaryGift string  `json:"anniversaryGift"`
			RankProgress    float64 `json:"rankProgress"`
		}{
			VipLevel:        vipLevelStr,
			NewVipLevel:     vipLevelStr,
			TotalBetAmount:  totalBetAmount,
			WeeklyTurnover:  weeklyTurnover,
			BirthdayBonus:   benefits.Birthdaybonus,
			BirthdayGift:    benefits.Birthdaygift,
			MonthlyGift:     benefits.Monthlygift,
			AnniversaryGift: benefits.Anniversarygift,
			RankProgress:    float64(ranking.Rankprogress.Exp),
		}

		jsonutils.Write(w, response, http.StatusOK)
	}
}
