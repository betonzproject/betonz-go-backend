package numericutils

import (
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

var ten = big.NewInt(10)

// Compares two numerics and returns:
//
//	-1 if n1 <  n2
//	 0 if n1 == n2
//	+1 if n1 >  n2
func Cmp(n1, n2 pgtype.Numeric) int {
	if n1.Exp == n2.Exp {
		return n1.Int.Cmp(n2.Int)
	}

	rn1, rn2 := RescalePair(n1, n2)

	return rn1.Int.Cmp(rn2.Int)
}

// RescalePair rescales two numerics to a common exponential value (min exp of both numerics)
func RescalePair(n1 pgtype.Numeric, n2 pgtype.Numeric) (pgtype.Numeric, pgtype.Numeric) {
	if n1.Exp < n2.Exp {
		return n1, rescale(n2, n1.Exp)
	} else if n1.Exp > n2.Exp {
		return rescale(n1, n2.Exp), n2
	}

	return n1, n2
}

func rescale(n pgtype.Numeric, exp int32) pgtype.Numeric {
	if n.Exp == exp {
		return n
	}

	// NOTE(vadim): must convert exps to float64 before - to prevent overflow
	diff := math.Abs(float64(exp) - float64(n.Exp))
	value := new(big.Int).Set(n.Int)

	expScale := new(big.Int).Exp(ten, big.NewInt(int64(diff)), nil)
	if exp > n.Exp {
		value = value.Quo(value, expScale)
	} else if exp < n.Exp {
		value = value.Mul(value, expScale)
	}

	return pgtype.Numeric{
		Int:              value,
		Exp:              exp,
		InfinityModifier: n.InfinityModifier,
		NaN:              n.NaN,
		Valid:            n.Valid,
	}
}
