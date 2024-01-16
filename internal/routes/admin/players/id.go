package players

import (
	"log"
	"net/http"
	"time"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Response struct {
	Player            db.GetPlayerInfoByIdRow              `json:"player"`
	RestrictionEvents []db.GetRestrictionEventsByUserIdRow `json:"restrictionEvents"`
	Banks             []db.Bank                            `json:"banks"`
	DailyTurnover     map[int32]int64                      `json:"dailyTurnover"`
	WeeklyTurnover    map[int32]int64                      `json:"weeklyTurnover"`
	MonthlyTurnover   map[int32]int64                      `json:"monthlyTurnover"`
	AlltimeTurnover   map[int32]int64                      `json:"alltimeTurnover"`
}

func GetPlayersById(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ManagePlayers) != nil {
			return
		}

		idParam := chi.URLParam(r, "id")
		id, err := utils.ParseUUID(idParam)
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		player, err := app.DB.GetPlayerInfoById(r.Context(), id)
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		restrictionEvents, err := app.DB.GetRestrictionEventsByUserId(r.Context(), id)
		if err != nil {
			log.Panicln("Error getting restriction events: " + err.Error())
		}

		banks, err := app.DB.GetBanksByUserId(r.Context(), id)
		if err != nil {
			log.Panicln("Error getting banks: " + err.Error())
		}

		location, _ := time.LoadLocation("Asia/Yangon")
		now := time.Now().In(location)

		dailyTurnover, err := app.DB.GetTurnoverByUserId(r.Context(), db.GetTurnoverByUserIdParams{
			ID:       id,
			FromDate: pgtype.Timestamptz{Time: timeutils.StartOfToday(), Valid: true},
			ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			log.Panicln("Error getting daily turnover: " + err.Error())
		}

		startOfLastWeek := time.Date(now.Year(), now.Month(), now.Day()-7, 0, 0, 0, 0, location)
		weeklyTurnover, err := app.DB.GetTurnoverByUserId(r.Context(), db.GetTurnoverByUserIdParams{
			ID:       id,
			FromDate: pgtype.Timestamptz{Time: startOfLastWeek, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			log.Panicln("Error getting weekly turnover: " + err.Error())
		}

		startOfLastMonth := time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, location)
		monthlyTurnover, err := app.DB.GetTurnoverByUserId(r.Context(), db.GetTurnoverByUserIdParams{
			ID:       id,
			FromDate: pgtype.Timestamptz{Time: startOfLastMonth, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			log.Panicln("Error getting monthly turnover: " + err.Error())
		}

		alltimeTurnover, err := app.DB.GetTurnoverByUserId(r.Context(), db.GetTurnoverByUserIdParams{
			ID:       id,
			FromDate: pgtype.Timestamptz{Valid: true},
			ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			log.Panicln("Error getting alltime turnover: " + err.Error())
		}

		jsonutils.Write(w, Response{
			Player:            player,
			RestrictionEvents: restrictionEvents,
			Banks:             banks,
			DailyTurnover:     toMap(dailyTurnover),
			WeeklyTurnover:    toMap(weeklyTurnover),
			MonthlyTurnover:   toMap(monthlyTurnover),
			AlltimeTurnover:   toMap(alltimeTurnover),
		}, http.StatusOK)
	}
}

func toMap(tos []db.GetTurnoverByUserIdRow) map[int32]int64 {
	turnoverMap := make(map[int32]int64)
	for _, turnover := range tos {
		turnoverMap[turnover.ProductCode] = turnover.Turnover
	}
	return turnoverMap
}
