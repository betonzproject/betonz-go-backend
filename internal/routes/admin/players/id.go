package players

import (
	"log"
	"net/http"
	"time"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/product"
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
	Turnover          map[string]int64                     `json:"turnover"`
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
		turnoverParam := r.URL.Query().Get("turnover")
		var turnover []db.GetTurnoverByUserIdRow

		switch turnoverParam {
		case "weekly":
			startOfLastWeek := time.Date(now.Year(), now.Month(), now.Day()-7, 0, 0, 0, 0, location)
			turnover, err = app.DB.GetTurnoverByUserId(r.Context(), db.GetTurnoverByUserIdParams{
				ID:       id,
				FromDate: pgtype.Timestamptz{Time: startOfLastWeek, Valid: true},
				ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			})
			if err != nil {
				log.Panicln("Error getting weekly turnover: " + err.Error())
			}

		case "monthly":
			startOfLastMonth := time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, location)
			turnover, err = app.DB.GetTurnoverByUserId(r.Context(), db.GetTurnoverByUserIdParams{
				ID:       id,
				FromDate: pgtype.Timestamptz{Time: startOfLastMonth, Valid: true},
				ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			})
			if err != nil {
				log.Panicln("Error getting monthly turnover: " + err.Error())
			}

		case "all-time":
			turnover, err = app.DB.GetTurnoverByUserId(r.Context(), db.GetTurnoverByUserIdParams{
				ID:       id,
				FromDate: pgtype.Timestamptz{Valid: true},
				ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			})
			if err != nil {
				log.Panicln("Error getting alltime turnover: " + err.Error())
			}

		default:
			turnover, err = app.DB.GetTurnoverByUserId(r.Context(), db.GetTurnoverByUserIdParams{
				ID:       id,
				FromDate: pgtype.Timestamptz{Time: timeutils.StartOfToday(), Valid: true},
				ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			})
			if err != nil {
				log.Panicln("Error getting daily turnover: " + err.Error())
			}
		}

		jsonutils.Write(w, Response{
			Player:            player,
			RestrictionEvents: restrictionEvents,
			Banks:             banks,
			Turnover:          toMap(turnover),
		}, http.StatusOK)
	}
}

func toMap(tos []db.GetTurnoverByUserIdRow) map[string]int64 {
	turnoverMap := make(map[string]int64)
	for _, p := range product.AllProducts {
		turnoverMap[p.String()] = 0
	}
	for _, turnover := range tos {
		turnoverMap[product.Product(int(turnover.ProductCode)).String()] = turnover.Turnover
	}
	return turnoverMap
}
