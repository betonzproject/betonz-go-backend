package promotion

import (
	"context"
	"log"
	"math/big"

	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetEligiblePromotions(q *db.Queries, ctx context.Context, userId pgtype.UUID) []db.PromotionType {
	promotions := make([]db.PromotionType, 0, 3)
	if isEligibleForInactiveBonus(q, ctx, userId) {
		promotions = append(promotions, db.PromotionTypeINACTIVEBONUS)
	}
	return promotions
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

func CalculateBonus(amount pgtype.Numeric, promotion db.PromotionType) pgtype.Numeric {
	switch promotion {
	case db.PromotionTypeINACTIVEBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(2), Exp: -1, Valid: true})
	}
	return numericutils.Zero
}

func CalculateTurnoverTarget(amount pgtype.Numeric, promotion db.PromotionType) pgtype.Numeric {
	switch promotion {
	case db.PromotionTypeINACTIVEBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(10), Valid: true})
	}
	return numericutils.Zero
}
