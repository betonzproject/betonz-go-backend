package promotion

import (
	"context"
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
	hasApprovedDepositRequestsWithin30Days, _ := q.HasApprovedDepositRequestsWithin30DaysByUserId(ctx, userId)
	hasActiveInactiveBonus, _ := q.HasActivePromotionByUserId(ctx, db.HasActivePromotionByUserIdParams{
		UserId:    userId,
		PromoCode: db.NullPromotionType{PromotionType: db.PromotionTypeINACTIVEBONUS, Valid: true},
	})
	hasPendingInactiveBonus, _ := q.HasPendingTransactionRequestsWithPromotion(ctx, db.HasPendingTransactionRequestsWithPromotionParams{
		UserId:    userId,
		Promotion: db.NullPromotionType{PromotionType: db.PromotionTypeINACTIVEBONUS, Valid: true},
	})

	return (!hasApprovedDepositRequestsWithin30Days &&
		!hasActiveInactiveBonus &&
		!hasPendingInactiveBonus)
}

func GetBonus(amount pgtype.Numeric, promotion db.PromotionType) pgtype.Numeric {
	switch promotion {
	case db.PromotionTypeINACTIVEBONUS:
		return numericutils.Mul(amount, pgtype.Numeric{Int: big.NewInt(2), Exp: -1, Valid: true})
	}
	return numericutils.Zero
}
