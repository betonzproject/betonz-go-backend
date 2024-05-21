package report

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/timeutils"

	"github.com/jackc/pgx/v5/pgtype"
)

type RetentionReport struct {
	TotalPlayerCount  int64 `json:"totalPlayers"`
	ActivePlayerCount int64 `json:"activeUserCount"`
}

type MonthlyReportRow struct {
	Date                                time.Time `json:"date"`
	WinLoss                             int64     `json:"winLoss"`
	WinLossIncreasePercentage           int       `json:"winLossPercentage"`
	DepositTotal                        int64     `json:"depositTotal"`
	DepositTotalIncreasePercentage      int       `json:"depositTotalIncreasePercentage"`
	DepositCount                        int64     `json:"depositCount"`
	DepositCountIncreasePercentage      int       `json:"depositCountIncreasePercentage"`
	WithdrawTotal                       int64     `json:"withdrawTotal"`
	WithdrawTotalIncreasePercentage     int       `json:"withdrawTotalIncreasePercentage"`
	WithdrawCount                       int64     `json:"withdrawCount"`
	WithdrawCountIncreasePercentage     int       `json:"withdrawCountIncreasePercentage"`
	BonusTotal                          int64     `json:"bonusTotal"`
	BonusTotalIncreasePercentage        int       `json:"bonusTotalIncreasePercentage"`
	ActivePlayerCount                   int64     `json:"activePlayerCount"`
	ActivePlayerCountIncreasePercentage int       `json:"activePlayerCountIncreasePercentage"`
}

type ReportResponse struct {
	RetentionReport RetentionReport             `json:"retentionReport"`
	MonthlyReport   []MonthlyReportRow          `json:"monthlyReport"`
	DailyReport     []db.GetDailyPerformanceRow `json:"dailyReport"`
}

func GetReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.Authenticate(app, w, r)
		if err != nil {
			return
		}

		if acl.Authorize(app, w, r, user.Role, acl.ViewReports) != nil {
			return
		}

		retentionFromParam := r.URL.Query().Get("retentionFrom")
		retentionToParam := r.URL.Query().Get("retentionTo")

		retentionFrom, err := timeutils.ParseDatetime(retentionFromParam)
		if err != nil {
			retentionFrom = timeutils.StartOfToday().AddDate(0, 0, -6)
		}
		retentionTo, err := timeutils.ParseDatetime(retentionToParam)
		if err != nil {
			retentionTo = timeutils.EndOfToday()
		}

		retentionReport := getRetentionReport(app, r.Context(), retentionFrom, retentionTo)

		monthlyReportDateParam := r.URL.Query().Get("monthlyReportDate")
		monthlyReportDate, err := timeutils.ParseDatetime(monthlyReportDateParam)
		if err != nil {
			monthlyReportDate = time.Now()
		}

		monthlyReport := getMonthlyReport(app, r.Context(), monthlyReportDate)

		dailyReportFromParam := r.URL.Query().Get("dailyReportFrom")
		dailyReportToParam := r.URL.Query().Get("dailyReportTo")

		dailyReportFrom, err := timeutils.ParseDatetime(dailyReportFromParam)
		if err != nil {
			dailyReportFrom = timeutils.StartOfThisMonth()
		}
		dailyReportTo, err := timeutils.ParseDatetime(dailyReportToParam)
		if err != nil {
			dailyReportTo = timeutils.EndOfThisMonth()
		}

		dailyReport, err := app.DB.GetDailyPerformance(r.Context(), db.GetDailyPerformanceParams{
			FromDate: pgtype.Timestamptz{Time: dailyReportFrom, Valid: true},
			ToDate:   pgtype.Timestamptz{Time: dailyReportTo, Valid: true},
		})
		if err != nil {
			log.Panicln("Can't get daily performance data: ", err)
		}

		jsonutils.Write(w, ReportResponse{
			RetentionReport: retentionReport,
			MonthlyReport:   monthlyReport,
			DailyReport:     dailyReport,
		}, http.StatusOK)
	}
}

func getRetentionReport(app *app.App, ctx context.Context, from, to time.Time) RetentionReport {
	totalPlayerCount, err := app.DB.GetNewPlayerCount(ctx, db.GetNewPlayerCountParams{
		FromDate: pgtype.Timestamptz{Valid: true},
		ToDate:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get total players: " + err.Error())
	}

	activeUserCount, err := app.DB.GetActivePlayerCount(ctx, db.GetActivePlayerCountParams{
		FromDate: pgtype.Timestamptz{Time: from, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: to, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get active player count: " + err.Error())
	}

	return RetentionReport{
		TotalPlayerCount:  totalPlayerCount,
		ActivePlayerCount: activeUserCount,
	}
}

func getMonthlyReport(app *app.App, ctx context.Context, month time.Time) []MonthlyReportRow {
	location, _ := time.LoadLocation("Asia/Yangon")
	monthStart := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, location)
	monthEnd := time.Date(month.Year(), month.Month()+1, 1, 0, 0, 0, -1, location)
	lastMonthStart := monthStart.AddDate(0, -1, 0)
	lastMonthEnd := monthEnd.AddDate(0, -1, 0)
	secondLastMonthStart := monthStart.AddDate(0, -2, 0)
	secondLastMonthEnd := monthEnd.AddDate(0, -2, 0)

	row1 := getMonthlyReportRow(app, ctx, monthStart, monthEnd)
	row2 := getMonthlyReportRow(app, ctx, lastMonthStart, lastMonthEnd)
	row3 := getMonthlyReportRow(app, ctx, secondLastMonthStart, secondLastMonthEnd)

	return []MonthlyReportRow{row1, row2, row3}
}

func getMonthlyReportRow(app *app.App, ctx context.Context, monthStart, monthEnd time.Time) MonthlyReportRow {
	lastMonthStart := monthStart.AddDate(0, -1, 0)
	lastMonthEnd := monthEnd.AddDate(0, -1, 0)

	winLossThisMonth, err := app.DB.GetTotalWinLoss(ctx, db.GetTotalWinLossParams{
		FromDate: pgtype.Timestamptz{Time: monthStart, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: monthEnd, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get total winloss this month: " + err.Error())
	}

	winLossLastMonth, err := app.DB.GetTotalWinLoss(ctx, db.GetTotalWinLossParams{
		FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get total winloss last month: " + err.Error())
	}

	winLossIncreasePercentage := getIncreasePercentage(winLossLastMonth, winLossThisMonth)

	depositsThisMonth, err := app.DB.GetTotalTransactionAmountAndCount(ctx, db.GetTotalTransactionAmountAndCountParams{
		Type:     db.TransactionTypeDEPOSIT,
		FromDate: pgtype.Timestamptz{Time: monthStart, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: monthEnd, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get deposit stats this month: " + err.Error())
	}

	depositsLastMonth, err := app.DB.GetTotalTransactionAmountAndCount(ctx, db.GetTotalTransactionAmountAndCountParams{
		Type:     db.TransactionTypeDEPOSIT,
		FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get deposit stats last month: " + err.Error())
	}

	depositTotalIncreasePercentage := getIncreasePercentage(depositsLastMonth.Total, depositsThisMonth.Total)
	depositCountIncreasePercentage := getIncreasePercentage(depositsLastMonth.Count, depositsThisMonth.Count)

	withdrawsThisMonth, err := app.DB.GetTotalTransactionAmountAndCount(ctx, db.GetTotalTransactionAmountAndCountParams{
		Type:     db.TransactionTypeWITHDRAW,
		FromDate: pgtype.Timestamptz{Time: monthStart, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: monthEnd, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get withdraw stats this month: " + err.Error())
	}

	withdrawsLastMonth, err := app.DB.GetTotalTransactionAmountAndCount(ctx, db.GetTotalTransactionAmountAndCountParams{
		Type:     db.TransactionTypeWITHDRAW,
		FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get withdraw stats last month: " + err.Error())
	}

	withdrawTotalIncreasePercentage := getIncreasePercentage(withdrawsLastMonth.Total, withdrawsThisMonth.Total)
	withdrawCountIncreasePercentage := getIncreasePercentage(withdrawsLastMonth.Count, withdrawsThisMonth.Count)

	bonusTotalIncreasePercentage := getIncreasePercentage(depositsLastMonth.BonusTotal, depositsThisMonth.BonusTotal)

	activePlayerCountThisMonth, err := app.DB.GetActivePlayerCount(ctx, db.GetActivePlayerCountParams{
		FromDate: pgtype.Timestamptz{Time: monthStart, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: monthEnd, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get active player count this month: " + err.Error())
	}

	activePlayerCountLastMonth, err := app.DB.GetActivePlayerCount(ctx, db.GetActivePlayerCountParams{
		FromDate: pgtype.Timestamptz{Time: lastMonthStart, Valid: true},
		ToDate:   pgtype.Timestamptz{Time: lastMonthEnd, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get active player count last month: " + err.Error())
	}

	activePlayerCountIncreasePercentage := getIncreasePercentage(activePlayerCountLastMonth, activePlayerCountThisMonth)

	return MonthlyReportRow{
		Date:                                monthStart,
		WinLoss:                             winLossThisMonth,
		WinLossIncreasePercentage:           winLossIncreasePercentage,
		DepositTotal:                        depositsThisMonth.Total,
		DepositTotalIncreasePercentage:      depositTotalIncreasePercentage,
		DepositCount:                        depositsThisMonth.Count,
		DepositCountIncreasePercentage:      depositCountIncreasePercentage,
		WithdrawTotal:                       withdrawsThisMonth.Total,
		WithdrawTotalIncreasePercentage:     withdrawTotalIncreasePercentage,
		WithdrawCount:                       withdrawsThisMonth.Count,
		WithdrawCountIncreasePercentage:     withdrawCountIncreasePercentage,
		BonusTotal:                          depositsThisMonth.BonusTotal,
		BonusTotalIncreasePercentage:        bonusTotalIncreasePercentage,
		ActivePlayerCount:                   activePlayerCountThisMonth,
		ActivePlayerCountIncreasePercentage: activePlayerCountIncreasePercentage,
	}
}

func getIncreasePercentage(before, after int64) int {
	var percentage int
	if before == 0 {
		if after > 0 {
			percentage = 100
		}
	} else {
		percentage = int(float64(after-before) / float64(before) * 100)
	}
	return percentage
}
