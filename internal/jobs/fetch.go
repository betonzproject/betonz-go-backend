package jobs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/etg"
	"github.com/doorman2137/betonz-go/internal/utils/transactionutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type EtgTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05"

func (ct *EtgTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}

func (ct *EtgTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
}

type FetchRequest struct {
	Op  string `json:"op"`
	Key int    `json:"key"`
}

type FetchReponseData struct {
	BetId      int32          `json:"betid"`
	RefId      string         `json:"refid"`
	Op         string         `json:"op"`
	Member     string         `json:"member"`
	Prodmember string         `json:"prodmember"`
	Prodcode   int32          `json:"prodcode"`
	Prodtype   int32          `json:"prodtype"`
	GameId     string         `json:"gameid"`
	BetDetail  string         `json:"betdetail"`
	Turnover   pgtype.Numeric `json:"turnover"`
	Bet        pgtype.Numeric `json:"bet"`
	Payout     pgtype.Numeric `json:"payout"`
	Winlose    pgtype.Numeric `json:"winlose"`
	Status     int32          `json:"status"`
	Starttime  EtgTime        `json:"starttime"`
	Matchtime  EtgTime        `json:"matchtime"`
	Endtime    EtgTime        `json:"endtime"`
	Settletime EtgTime        `json:"settletime"`
	Progshare  pgtype.Numeric `json:"progshare"`
	Progwin    pgtype.Numeric `json:"progwin"`
	Comm       pgtype.Numeric `json:"comm"`
}

type FetchResponse struct {
	Err     int                `json:"err"`
	Desc    string             `json:"desc"`
	NextKey int                `json:"nextkey"`
	Data    []FetchReponseData `json:"data"`
}

type MarkRequest struct {
	Op   string  `json:"op"`
	Mark []int32 `json:"mark"`
}

type MarkResponse struct {
	Err  int    `json:"err"`
	Desc string `json:"desc"`
}

func FetchBets(app *app.App, key int) {
	ctx := context.Background()
	now := time.Now()

	payload := FetchRequest{
		Op:  os.Getenv("ETG_OPERATOR_CODE"),
		Key: key,
	}
	var fetchResponse FetchResponse
	err := etg.Post("/fetch", payload, &fetchResponse)
	if err != nil {
		log.Println("Can't fetch bets: " + err.Error())
		return
	}

	if fetchResponse.Err != etg.Success {
		log.Printf("Can't fetch bet records with key %d: %s\nEndpoint: %s\nPayload: %+v\n", key, err, "/fetch", payload)
		return
	}

	if len(fetchResponse.Data) == 0 {
		log.Println("No records to fetch")
		return
	}

	tx, qtx := transactionutils.Begin(app, ctx)
	defer tx.Rollback(ctx)

	idsToMark := make([]int32, 0, 100)
	for _, bet := range fetchResponse.Data {
		err := qtx.UpsertBet(ctx, db.UpsertBetParams{
			ID:               bet.BetId,
			RefId:            bet.RefId,
			EtgUsername:      bet.Member,
			ProviderUsername: bet.Prodmember,
			ProductCode:      bet.Prodcode,
			ProductType:      bet.Prodtype,
			GameId:           pgtype.Text{String: bet.GameId, Valid: true},
			Details:          bet.BetDetail,
			WinLoss:          bet.Winlose,
			Turnover:         bet.Turnover,
			Bet:              bet.Bet,
			Payout:           bet.Payout,
			Status:           bet.Status,
			StartTime:        pgtype.Timestamptz{Time: bet.Starttime.Time, Valid: true},
			MatchTime:        pgtype.Timestamptz{Time: bet.Matchtime.Time, Valid: true},
			EndTime:          pgtype.Timestamptz{Time: bet.Endtime.Time, Valid: true},
			SettleTime:       pgtype.Timestamptz{Time: bet.Settletime.Time, Valid: true},
			ProgShare:        bet.Progshare,
			ProgWin:          bet.Progwin,
			Commission:       bet.Comm,
		})
		if err != nil {
			log.Println("Can't upsert bet: ", err.Error())
			return
		}

		if bet.Member != "demo001" {
			idsToMark = append(idsToMark, bet.BetId)
		}
	}

	if os.Getenv("ENVIRONMENT") != "development" {
		markPayload := MarkRequest{
			Op:   os.Getenv("ETG_OPERATOR_CODE"),
			Mark: idsToMark,
		}

		var markResponse MarkResponse
		err := etg.Post("/mark", markPayload, &markResponse)
		if err != nil || markResponse.Err != etg.Success {
			log.Println("Can't mark bets: " + err.Error())
			return
		}
	}

	d := time.Since(now)
	log.Printf("Fetched %d bets, marked %d bets in %s", len(fetchResponse.Data), len(idsToMark), d)

	tx.Commit(ctx)
}
