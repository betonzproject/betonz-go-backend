package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils/timeutils"
)

type DailyReward struct {
	Day    int                  `json:"day"`
	Reward db.InventoryItemType `json:"reward"`
	Count  int                  `json:"count"`
	Amount int                  `json:"amount"`
}

type GetRewardResponse struct {
	ClaimedRewardCount int           `json:"claimedRewardCount"`
	UpcomingReward     int           `json:"upcomingReward"`
	Rewards            []DailyReward `json:"rewards"`
}

type RewardEventData struct {
	LastClaimCount int32 `json:"lastClaimCount"`
}

// This function will be return an array that include reward object
// @params maxDay is the count of object in the array
// if maxDay is 29, the rewards object in array will be 29 object
func GetRewards(maxDay int) ([]DailyReward, error) {
	rewards := map[int]DailyReward{
		1:  {Day: 1, Reward: db.InventoryItemTypeTOKENA, Count: 3},
		2:  {Day: 2, Reward: db.InventoryItemTypeTOKENB, Count: 1},
		3:  {Day: 3, Reward: db.InventoryItemTypeBETONPOINT, Count: 2000},
		4:  {Day: 4, Reward: db.InventoryItemTypeBONUS, Amount: 2000},
		5:  {Day: 5, Reward: db.InventoryItemTypeREDPACK, Count: 1},
		6:  {Day: 6, Reward: db.InventoryItemTypeTOKENA, Count: 3},
		7:  {Day: 7, Reward: db.InventoryItemTypeTOKENB, Count: 1},
		8:  {Day: 8, Reward: db.InventoryItemTypeBETONPOINT, Count: 2000},
		9:  {Day: 9, Reward: db.InventoryItemTypeBONUS, Amount: 3000},
		10: {Day: 10, Reward: db.InventoryItemTypeREDPACK, Count: 1},
		11: {Day: 11, Reward: db.InventoryItemTypeTOKENA, Count: 3},
		12: {Day: 12, Reward: db.InventoryItemTypeTOKENB, Count: 1},
		13: {Day: 13, Reward: db.InventoryItemTypeBETONPOINT, Count: 2000},
		14: {Day: 14, Reward: db.InventoryItemTypeBONUS, Amount: 3000},
		15: {Day: 15, Reward: db.InventoryItemTypeREDPACK, Count: 1},
		16: {Day: 16, Reward: db.InventoryItemTypeTOKENA, Count: 6},
		17: {Day: 17, Reward: db.InventoryItemTypeTOKENB, Count: 2},
		18: {Day: 18, Reward: db.InventoryItemTypeBETONPOINT, Count: 4000},
		19: {Day: 19, Reward: db.InventoryItemTypeBONUS, Amount: 4000},
		20: {Day: 20, Reward: db.InventoryItemTypeREDPACK, Count: 1},
		21: {Day: 21, Reward: db.InventoryItemTypeTOKENA, Count: 6},
		22: {Day: 22, Reward: db.InventoryItemTypeTOKENB, Count: 2},
		23: {Day: 23, Reward: db.InventoryItemTypeBETONPOINT, Count: 4000},
		24: {Day: 24, Reward: db.InventoryItemTypeBONUS, Amount: 6000},
		25: {Day: 25, Reward: db.InventoryItemTypeROYALREDPACK, Count: 1},
		26: {Day: 26, Reward: db.InventoryItemTypeTOKENA, Count: 6},
		27: {Day: 27, Reward: db.InventoryItemTypeTOKENB, Count: 2},
		28: {Day: 28, Reward: db.InventoryItemTypeBETONPOINT, Count: 4000},
		29: {Day: 29, Reward: db.InventoryItemTypeRAFFLETICKET, Count: 1},
		30: {Day: 30, Reward: db.InventoryItemTypeROYALREDPACK, Count: 1},
		31: {Day: 31, Reward: db.InventoryItemTypeBONUS, Amount: 10000},
	}

	rewardList := make([]DailyReward, 0, maxDay)
	for i := 1; i <= maxDay; i++ {
		reward, ok := rewards[i]
		if !ok {
			return nil, errors.New("no reward for this day")
		}
		rewardList = append(rewardList, reward)
	}

	return rewardList, nil
}

// This function will be return an object that include reward information
func CheckAvaliableReward(q *db.Queries, r *http.Request, user db.GetUserByIdRow) GetRewardResponse {
	dailyRewards, _ := GetRewards(timeutils.DaysInMonth())

	lastRewardById, err := q.GetLastRewardClaimedById(r.Context(), user.ID)
	if err != nil {
		return GetRewardResponse{
			ClaimedRewardCount: 0,
			UpcomingReward:     1,
			Rewards:            dailyRewards,
		}
	}

	dataJSON, err := json.Marshal(lastRewardById.Data)
	if err != nil {
		log.Panicln("Cannot parse data from claimed reward event: ", err.Error())
	}

	var lastRewardData RewardEventData
	if err := json.Unmarshal(dataJSON, &lastRewardData); err != nil {
		log.Panicln("Cannot parse lastClaim reward count: ", err.Error())
	}

	if lastRewardData.LastClaimCount >= int32(timeutils.DaysInMonth()) || lastRewardData.LastClaimCount == 0 {
		return GetRewardResponse{
			ClaimedRewardCount: 0,
			UpcomingReward:     1,
			Rewards:            dailyRewards,
		}
	} else {
		return GetRewardResponse{
			ClaimedRewardCount: int(lastRewardData.LastClaimCount),
			UpcomingReward:     int(lastRewardData.LastClaimCount) + 1,
			Rewards:            dailyRewards,
		}
	}
}
