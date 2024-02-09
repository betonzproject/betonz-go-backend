package report

import (
	"log"
	"net/http"

	"github.com/doorman2137/betonz-go/internal/acl"
	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/auth"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/doorman2137/betonz-go/internal/utils/timeutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type OverviewReport struct {
	DepositAmountToday                  int64 `json:"depositAmountToday"`
	DepositCountToday                   int64 `json:"depositCountToday"`
	DepositAmountYesterday              int64 `json:"depositAmountYesterday"`
	DepositCountYesterday               int64 `json:"depositCountYesterday"`
	DepositAmountThisMonth              int64 `json:"depositAmountThisMonth"`
	DepositCountThisMonth               int64 `json:"depositCountThisMonth"`
	DepositAmountLastMonth              int64 `json:"depositAmountLastMonth"`
	DepositCountLastMonth               int64 `json:"depositCountLastMonth"`
	WithdrawAmountToday                 int64 `json:"withdrawAmountToday"`
	WithdrawCountToday                  int64 `json:"withdrawCountToday"`
	WithdrawAmountYesterday             int64 `json:"withdrawAmountYesterday"`
	WithdrawCountYesterday              int64 `json:"withdrawCountYesterday"`
	WithdrawAmountThisMonth             int64 `json:"withdrawAmountThisMonth"`
	WithdrawCountThisMonth              int64 `json:"withdrawCountThisMonth"`
	WithdrawAmountLastMonth             int64 `json:"withdrawAmountLastMonth"`
	WithdrawCountLastMonth              int64 `json:"withdrawCountLastMonth"`
	NewPlayersThisMonth                 int64 `json:"newPlayersThisMonth"`
	NewPlayersLastMonth                 int64 `json:"newPlayersLastMonth"`
	NewPlayersThisYear                  int64 `json:"newPlayersThisYear"`
	NewPlayersLastYear                  int64 `json:"newPlayersLastYear"`
	PlayersWithTransactionsThisMonth    int64 `json:"playersWithTransactionsThisMonth"`
	PlayersWithTransactionsLastMonth    int64 `json:"playersWithTransactionsLastMonth"`
	NewPlayersWithTransactionsThisMonth int64 `json:"newPlayersWithTransactionsThisMonth"`
	NewPlayersWithTransactionsLastMonth int64 `json:"newPlayersWithTransactionsLastMonth"`
}

func GetOverview(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewReports) != nil {
			return
		}

		todayStart := timeutils.StartOfToday()
		todayEnd := timeutils.EndOfToday()
		yesterdayStart := todayStart.AddDate(0, 0, -1)
		yesterdayEnd := todayStart.AddDate(0, 0, -1)
		thisMonthStart := timeutils.StartOfThisMonth()
		thisMonthEnd := timeutils.EndOfThisMonth()
		lastMonthStart := thisMonthStart.AddDate(0, -1, 0)
		lastMonthEnd := thisMonthEnd.AddDate(0, -1, 0)
		thisYearStart := timeutils.StartOfThisYear()
		thisYearEnd := timeutils.EndOfThisYear()
		lastYearStart := thisYearStart.AddDate(-1, 0, 0)
		lastYearEnd := thisYearEnd.AddDate(-1, 0, 0)

		depositsToday, err := app.DB.GetTotalTransactionAmountAndCount(r.Context(), db.GetTotalTransactionAmountAndCountParams{
			Type:     db.TransactionTypeDEPOSIT,
			FromDate: pgtype.Timestamptz{Time: todayStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: todayEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get deposit stats today: " + err.Error())
		}

		depositsYesterday, err := app.DB.GetTotalTransactionAmountAndCount(r.Context(), db.GetTotalTransactionAmountAndCountParams{
			Type:     db.TransactionTypeDEPOSIT,
			FromDate: pgtype.Timestamptz{Time: yesterdayStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: yesterdayEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get deposit stats yesterday: " + err.Error())
		}

		depositsThisMonth, err := app.DB.GetTotalTransactionAmountAndCount(r.Context(), db.GetTotalTransactionAmountAndCountParams{
			Type:     db.TransactionTypeDEPOSIT,
			FromDate: pgtype.Timestamptz{Time: thisMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: thisMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get deposit stats this month: " + err.Error())
		}

		depositsLastMonth, err := app.DB.GetTotalTransactionAmountAndCount(r.Context(), db.GetTotalTransactionAmountAndCountParams{
			Type:     db.TransactionTypeDEPOSIT,
			FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get deposit stats last month: " + err.Error())
		}

		withdrawalsToday, err := app.DB.GetTotalTransactionAmountAndCount(r.Context(), db.GetTotalTransactionAmountAndCountParams{
			Type:     db.TransactionTypeWITHDRAW,
			FromDate: pgtype.Timestamptz{Time: todayStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: todayEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get withdrawal stats today: " + err.Error())
		}

		withdrawalsYesterday, err := app.DB.GetTotalTransactionAmountAndCount(r.Context(), db.GetTotalTransactionAmountAndCountParams{
			Type:     db.TransactionTypeWITHDRAW,
			FromDate: pgtype.Timestamptz{Time: yesterdayStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: yesterdayEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get withdrawal stats today: " + err.Error())
		}

		withdrawalsThisMonth, err := app.DB.GetTotalTransactionAmountAndCount(r.Context(), db.GetTotalTransactionAmountAndCountParams{
			Type:     db.TransactionTypeWITHDRAW,
			FromDate: pgtype.Timestamptz{Time: thisMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: thisMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get withdrawal stats this month: " + err.Error())
		}

		withdrawalsLastMonth, err := app.DB.GetTotalTransactionAmountAndCount(r.Context(), db.GetTotalTransactionAmountAndCountParams{
			Type:     db.TransactionTypeWITHDRAW,
			FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get withdrawal stats last month: " + err.Error())
		}

		newPlayerCountThisMonth, err := app.DB.GetNewPlayerCount(r.Context(), db.GetNewPlayerCountParams{
			FromDate: pgtype.Timestamptz{Time: thisMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: thisMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get new player count this month: " + err.Error())
		}

		newPlayerCountLastMonth, err := app.DB.GetNewPlayerCount(r.Context(), db.GetNewPlayerCountParams{
			FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get new player count last month: " + err.Error())
		}

		newPlayerCountThisYear, err := app.DB.GetNewPlayerCount(r.Context(), db.GetNewPlayerCountParams{
			FromDate: pgtype.Timestamptz{Time: thisYearStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: thisYearEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get new player count this year: " + err.Error())
		}

		newPlayerCountLastYear, err := app.DB.GetNewPlayerCount(r.Context(), db.GetNewPlayerCountParams{
			FromDate: pgtype.Timestamptz{Time: lastYearStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: lastYearEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get new player count last year: " + err.Error())
		}

		playerWithTransactionsCountThisMonth, err := app.DB.GetPlayerWithTransactionsCount(r.Context(), db.GetPlayerWithTransactionsCountParams{
			FromDate: pgtype.Timestamptz{Time: thisMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: thisMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get player with transactions count this month: " + err.Error())
		}

		playerWithTransactionsCountLastMonth, err := app.DB.GetPlayerWithTransactionsCount(r.Context(), db.GetPlayerWithTransactionsCountParams{
			FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get player with transactions count last month: " + err.Error())
		}

		newPlayerWithTransactionsCountThisMonth, err := app.DB.GetPlayerWithTransactionsCount(r.Context(), db.GetPlayerWithTransactionsCountParams{
			FromDate: pgtype.Timestamptz{Time: thisMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: thisMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get new players with transactions count this month: " + err.Error())
		}

		newPlayerWithTransactionsCountLastMonth, err := app.DB.GetPlayerWithTransactionsCount(r.Context(), db.GetPlayerWithTransactionsCountParams{
			FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get new players with transactions count last month: " + err.Error())
		}

		jsonutils.Write(w, OverviewReport{
			DepositAmountToday:                  depositsToday.Total,
			DepositCountToday:                   depositsToday.Count,
			DepositAmountYesterday:              depositsYesterday.Total,
			DepositCountYesterday:               depositsYesterday.Count,
			DepositAmountThisMonth:              depositsThisMonth.Total,
			DepositCountThisMonth:               depositsThisMonth.Count,
			DepositAmountLastMonth:              depositsLastMonth.Total,
			DepositCountLastMonth:               depositsLastMonth.Count,
			WithdrawAmountToday:                 withdrawalsToday.Total,
			WithdrawCountToday:                  withdrawalsToday.Count,
			WithdrawAmountYesterday:             withdrawalsYesterday.Total,
			WithdrawCountYesterday:              withdrawalsYesterday.Count,
			WithdrawAmountThisMonth:             withdrawalsThisMonth.Total,
			WithdrawCountThisMonth:              withdrawalsThisMonth.Count,
			WithdrawAmountLastMonth:             withdrawalsLastMonth.Total,
			WithdrawCountLastMonth:              withdrawalsLastMonth.Count,
			NewPlayersThisMonth:                 newPlayerCountThisMonth,
			NewPlayersLastMonth:                 newPlayerCountLastMonth,
			NewPlayersThisYear:                  newPlayerCountThisYear,
			NewPlayersLastYear:                  newPlayerCountLastYear,
			PlayersWithTransactionsThisMonth:    playerWithTransactionsCountThisMonth,
			PlayersWithTransactionsLastMonth:    playerWithTransactionsCountLastMonth,
			NewPlayersWithTransactionsThisMonth: newPlayerWithTransactionsCountThisMonth,
			NewPlayersWithTransactionsLastMonth: newPlayerWithTransactionsCountLastMonth,
		}, http.StatusOK)
	}
}
