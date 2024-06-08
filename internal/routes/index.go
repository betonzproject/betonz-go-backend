package routes

import (
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/BetOnz-Company/betonz-go/internal/acl"
	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/auth"
	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/product"
	"github.com/BetOnz-Company/betonz-go/internal/utils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/jsonutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/numericutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/timeutils"
	"github.com/BetOnz-Company/betonz-go/internal/utils/transactionutils"

	"github.com/jackc/pgx/v5/pgtype"
)

type Response struct {
	User                    *db.GetUserByIdRow      `json:"user"`
	UnreadNotificationCount int64                   `json:"unreadNotificationCount"`
	Permissons              []acl.Permission        `json:"permissions"`
	ExpTarget               int64                   `json:"expTarget"`
	Reward                  utils.GetRewardResponse `json:"reward"`
	IsRewardClaimed         bool                    `json:"isRewardClaimed"`
}

func GetIndex(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isRewardClaimed := ClaimReward(app, w, r)

		user, err := auth.GetUser(app, w, r)
		if err != nil {
			jsonutils.Write(w, Response{}, http.StatusOK)
			return
		}

		userLevel, _ := user.Level.Int64Value()
		userExp, _ := user.Exp.Int64Value()

		nextLevelExp := utils.AllTargets[userLevel.Int64-1]

		canIncreaseLevel := utils.ExpTarget(userExp.Int64) >= nextLevelExp

		if canIncreaseLevel && userLevel.Int64 != 80 {
			err = app.DB.IncreaseUserLevelAndExp(r.Context(), db.IncreaseUserLevelAndExpParams{
				ID:  user.ID,
				Exp: pgtype.Numeric{Int: big.NewInt(int64(utils.ExpTarget(userExp.Int64) - nextLevelExp)), Valid: true},
			})
			if err != nil {
				log.Panicln("Error updating user's level: ", err.Error())
			}
		}

		var unreadNotificationCount int64
		if acl.IsAuthorized(user.Role, acl.ViewNotifications) {
			unreadNotificationCount, err = app.DB.GetUnreadNotificationCountByUserId(r.Context(), user.ID)
			if err != nil {
				log.Panicln("Can't get notification count: ", err)
			}
		}

		event, err := app.DB.GetActiveEventTodayByUserId(r.Context(), user.ID)
		if err != nil {
			err = utils.LogEvent(app.DB, r, user.ID, db.EventTypeACTIVE, db.EventResultSUCCESS, "", nil)
			if err != nil {
				log.Panicln("Can't log event: " + err.Error())
			}
		} else {
			err = app.DB.UpdateEvent(r.Context(), event.ID)
			if err != nil {
				log.Panicln("Can't update event: " + err.Error())
			}
		}

		reward := utils.CheckAvaliableReward(app.DB, r, user)

		if canIncreaseLevel && userLevel.Int64 != 80 {
			jsonutils.Write(w, Response{User: &user, UnreadNotificationCount: unreadNotificationCount, Permissons: acl.Acl[user.Role], ExpTarget: int64(utils.ExpTarget(userExp.Int64) - nextLevelExp), Reward: reward, IsRewardClaimed: isRewardClaimed}, http.StatusOK)
		} else {
			jsonutils.Write(w, Response{User: &user, UnreadNotificationCount: unreadNotificationCount, Permissons: acl.Acl[user.Role], ExpTarget: int64(nextLevelExp), Reward: reward, IsRewardClaimed: isRewardClaimed}, http.StatusOK)
		}
	}
}

func ClaimReward(app *app.App, w http.ResponseWriter, r *http.Request) bool {
	user, err := auth.GetUser(app, w, r)
	if err != nil {
		return false
	}

	// Check that the reward was already claimed within 24 hrs
	isAlreadyClaimed, _ := app.DB.HasRecentClaimedRewardByUserId(r.Context(), user.ID)
	if isAlreadyClaimed {
		return false
	}

	tx, qtx := transactionutils.Begin(app, r.Context())
	defer tx.Rollback(r.Context())

	lastRewardById, err := app.DB.GetLastRewardClaimedById(r.Context(), user.ID)
	// If user login for the first time, he will receive the first reward
	if err != nil {
		err = qtx.AddItemToInventory(r.Context(), db.AddItemToInventoryParams{
			UserId: user.ID,
			Item:   db.InventoryItemTypeTOKENA,
			Count:  pgtype.Numeric{Int: big.NewInt(3), Valid: true},
		})
		if err != nil {
			log.Panicln("Error inserting to Inventory table: ", err.Error())
		}

		err = utils.LogEvent(qtx, r, user.ID, db.EventTypeREWARDCLAIM, db.EventResultSUCCESS, "", map[string]any{
			"lastClaimCount": 1,
		})
		if err != nil {
			log.Panicln("Error creating event", err.Error())
		}

		tx.Commit(r.Context())
		return true
	}

	dataJSON, err := json.Marshal(lastRewardById.Data)
	if err != nil {
		return false
	}

	var lastRewardData utils.RewardEventData
	if err := json.Unmarshal(dataJSON, &lastRewardData); err != nil {
		log.Panicln("Cannot parse lastClaim reward count", err.Error())
	}

	// Reset Reward
	if lastRewardById.CreatedAt.Time.Month() != time.Now().Month() || lastRewardData.LastClaimCount >= int32(timeutils.DaysInMonth()) {
		err = utils.LogEvent(qtx, r, user.ID, db.EventTypeREWARDCLAIM, db.EventResultSUCCESS, "", map[string]any{
			"lastClaimCount": 0,
		})
		if err != nil {
			log.Panicln("Error creating event: ", err.Error())
		}
		tx.Commit(r.Context())
		return false
	}

	dailyRewards, _ := utils.GetRewards(timeutils.DaysInMonth())

	switch dailyRewards[lastRewardData.LastClaimCount].Reward {
	case db.InventoryItemTypeBONUS:
		err = qtx.DepositUserMainWallet(r.Context(), db.DepositUserMainWalletParams{
			ID: user.ID,
			Amount: pgtype.Numeric{
				Int:   big.NewInt(int64(dailyRewards[lastRewardData.LastClaimCount].Amount)),
				Valid: true,
			},
		})
		if err != nil {
			log.Panicln("Error depositing user's Main Wallet", err.Error())
		}
	case db.InventoryItemTypeBETONPOINT:
		// Add Beton Point
		err = qtx.AddUserBetonPoint(r.Context(), db.AddUserBetonPointParams{
			ID: user.ID,
			Bp: pgtype.Numeric{Int: big.NewInt(int64(dailyRewards[lastRewardData.LastClaimCount].Count)), Valid: true},
		})
		if err != nil {
			log.Panicln("Error adding beton point to user: ", err.Error())
		}
	default:
		// Add item to inventory
		err = qtx.AddItemToInventory(r.Context(), db.AddItemToInventoryParams{
			UserId: user.ID,
			Item:   dailyRewards[lastRewardData.LastClaimCount].Reward,
			Count:  pgtype.Numeric{Int: big.NewInt(int64(dailyRewards[lastRewardData.LastClaimCount].Count)), Valid: true},
		})
		if err != nil {
			log.Panicln("Error inserting to Inventory table: ", err.Error())
		}
	}

	err = utils.LogEvent(qtx, r, user.ID, db.EventTypeREWARDCLAIM, db.EventResultSUCCESS, "", map[string]any{
		"lastClaimCount": lastRewardData.LastClaimCount + 1,
	})
	if err != nil {
		log.Panicln("Error creating event", err.Error())
	}

	// Add transaction log event when the reward is bonus type
	if dailyRewards[lastRewardData.LastClaimCount].Reward == db.InventoryItemTypeBONUS {
		err = qtx.CreateTransactionRequest(r.Context(), db.CreateTransactionRequestParams{
			UserId: user.ID,
			Amount: pgtype.Numeric{
				Int:   big.NewInt(int64(dailyRewards[lastRewardData.LastClaimCount].Amount)),
				Valid: dailyRewards[lastRewardData.LastClaimCount].Reward == db.InventoryItemTypeBONUS,
			},
			DepositToWallet: pgtype.Int4{Int32: int32(product.MainWallet), Valid: true},
			Type:            db.TransactionTypeDEPOSIT,
			ReceiptPath:     pgtype.Text{Valid: true},
			Bonus:           numericutils.Zero,
			Status:          db.TransactionStatusAPPROVED,
			Remarks:         pgtype.Text{String: "Daily Reward", Valid: true}})
		if err != nil {
			log.Panicln("Error creating transaction request: ", err.Error())
		}
	}

	tx.Commit(r.Context())
	return true
}
