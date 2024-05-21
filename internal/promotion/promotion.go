package promotion

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/BetOnz-Company/betonz-go/internal/db"
	"github.com/BetOnz-Company/betonz-go/internal/utils/numericutils"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetEligiblePromotions(q *db.Queries, ctx context.Context, userId pgtype.UUID) ([]db.PromotionType, pgtype.Numeric, pgtype.Numeric) {
	promotions := make([]db.PromotionType, 0, 3)
	if isEligibleForInactiveBonus(q, ctx, userId) {
		promotions = append(promotions, db.PromotionTypeINACTIVEBONUS)
	}
	fivePercentBonusRemaining := getUnlimitedBonusRemaining(q, ctx, userId, db.PromotionTypeFIVEPERCENTUNLIMITEDBONUS)
	if numericutils.IsPositive(fivePercentBonusRemaining) {
		promotions = append(promotions, db.PromotionTypeFIVEPERCENTUNLIMITEDBONUS)
	}
	tenPercentBonusRemaining := getUnlimitedBonusRemaining(q, ctx, userId, db.PromotionTypeTENPERCENTUNLIMITEDBONUS)
	if numericutils.IsPositive(tenPercentBonusRemaining) {
		promotions = append(promotions, db.PromotionTypeTENPERCENTUNLIMITEDBONUS)
	}
	return promotions, fivePercentBonusRemaining, tenPercentBonusRemaining
}

func isEligibleForInactiveBonus(q *db.Queries, ctx context.Context, userId pgtype.UUID) bool {
	hasApprovedDepositRequestsWithin30Days, err := q.HasApprovedDepositRequestsWithin30DaysByUserId(ctx, userId)
	if err != nil {
		log.Panicln("Can't get approved deposits within 30 days: " + err.Error())
	}

	hasActiveInactiveBonus, err := q.HasActivePromotionByUserId(ctx, db.HasActivePromotionByUserIdParams{
		UserId:    userId,
		Promotion: db.NullPromotionType{PromotionType: db.PromotionTypeINACTIVEBONUS, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get active promotions: " + err.Error())
	}

	hasPendingInactiveBonus, err := q.HasPendingTransactionRequestsWithPromotion(ctx, db.HasPendingTransactionRequestsWithPromotionParams{
		UserId:    userId,
		Promotion: db.NullPromotionType{PromotionType: db.PromotionTypeINACTIVEBONUS, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get pending inactive bonuses: " + err.Error())
	}

	return (!hasApprovedDepositRequestsWithin30Days &&
		!hasActiveInactiveBonus &&
		!hasPendingInactiveBonus)
}

func getUnlimitedBonusRemaining(q *db.Queries, ctx context.Context, userId pgtype.UUID, promotion db.PromotionType) pgtype.Numeric {
	// Since this bonus resets at 8PM Myanmar time daily, create a special zone that is 4 hours ahead of Myanmar time (GMT+1030)
	location := time.FixedZone("Unlimited_Bonus_Zone", 37800) // 10.5*3600
	now := time.Now().In(location)
	fromDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	toDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, -1, location)

	remaining, err := q.GetBonusRemaining(ctx, db.GetBonusRemainingParams{
		UserId:    userId,
		Promotion: db.NullPromotionType{PromotionType: promotion, Valid: true},
		Limit:     pgtype.Numeric{Int: big.NewInt(1000000), Valid: true},
		FromDate:  pgtype.Timestamptz{Time: fromDate, Valid: true},
		ToDate:    pgtype.Timestamptz{Time: toDate, Valid: true},
	})
	if err != nil {
		log.Panicln("Can't get bonus remaining: " + err.Error())
	}

	return remaining
}

func CalculateBonus(amount pgtype.Numeric, promotion db.PromotionType) pgtype.Numeric {
	switch promotion {
	case db.PromotionTypeINACTIVEBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(20), Exp: -2, Valid: true})
	case db.PromotionTypeFIVEPERCENTUNLIMITEDBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(5), Exp: -2, Valid: true})
	case db.PromotionTypeTENPERCENTUNLIMITEDBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(10), Exp: -2, Valid: true})
	}
	return numericutils.Zero
}

func CalculateTurnoverTarget(amount pgtype.Numeric, promotion db.PromotionType) pgtype.Numeric {
	switch promotion {
	case db.PromotionTypeINACTIVEBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(10), Valid: true})
	case db.PromotionTypeFIVEPERCENTUNLIMITEDBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(2), Valid: true})
	case db.PromotionTypeTENPERCENTUNLIMITEDBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(6), Valid: true})
	}
	return numericutils.Zero
}
