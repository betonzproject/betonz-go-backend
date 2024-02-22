package numericutils

import (
	"log"
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

var ten = big.NewInt(10)

var Zero = pgtype.Numeric{Int: big.NewInt(0), Valid: true}
var NaN = pgtype.Numeric{NaN: true, Valid: true}

// Add returns the sum of all n.
// If any of the n is a NaN, then NaN is retruned.
//
// Panics if any n is not valid.
func Add(n ...pgtype.Numeric) pgtype.Numeric {
	if len(n) == 0 {
		return pgtype.Numeric{}
	}
	sum := n[0]
	for i := 1; i < len(n); i++ {
		sum = add(sum, n[i])
	}
	return sum
}

// add returns n1 + n2.
// If either n1 or n2 is NaN, NaN is returned.
//
// Panics if either n1 or n2 is not valid.
func add(n1, n2 pgtype.Numeric) pgtype.Numeric {
	ensureValid(n1)
	ensureValid(n2)

	if n1.NaN || n2.NaN {
		return NaN
	}

	rn1, rn2 := RescalePair(n1, n2)
	rn3 := new(big.Int).Add(rn1.Int, rn2.Int)
	return pgtype.Numeric{
		Int:              rn3,
		Exp:              rn1.Exp,
		InfinityModifier: n1.InfinityModifier,
		NaN:              n1.InfinityModifier-n2.InfinityModifier == 2 || n1.InfinityModifier-n2.InfinityModifier == -2,
		Valid:            true,
	}
}

// Add returns n1 + n2.
// If either n1 or n2 is NaN, NaN is returned.
//
// Panics if either n1 or n2 is not valid.
func Sub(n1, n2 pgtype.Numeric) pgtype.Numeric {
	ensureValid(n1)
	ensureValid(n2)

	if n1.NaN || n2.NaN {
		return NaN
	}

	rn1, rn2 := RescalePair(n1, n2)
	rn3 := new(big.Int).Sub(rn1.Int, rn2.Int)
	return pgtype.Numeric{
		Int:              rn3,
		Exp:              rn1.Exp,
		InfinityModifier: n1.InfinityModifier,
		NaN:              n1.InfinityModifier-n2.InfinityModifier == 2 || n1.InfinityModifier-n2.InfinityModifier == -2,
		Valid:            true,
	}
}

// Mul returns n1 * n2.
// If either n1 or n2 is NaN, NaN is returned.
//
// Panics if either n1 or n2 is not valid.
func Mul(n1, n2 pgtype.Numeric) pgtype.Numeric {
	ensureValid(n1)
	ensureValid(n2)

	expInt64 := int64(n1.Exp) + int64(n2.Exp)
	if expInt64 > math.MaxInt32 || expInt64 < math.MinInt32 {
		// NOTE(vadim): better to panic than give incorrect results, as
		// Decimals are usually used for money
		log.Panicf("exponent %v overflows an int32!\n", expInt64)
	}

	n3 := new(big.Int).Mul(n1.Int, n2.Int)
	return pgtype.Numeric{
		Int:              n3,
		Exp:              int32(expInt64),
		InfinityModifier: n1.InfinityModifier,
		NaN:              n1.InfinityModifier-n2.InfinityModifier == 2 || n1.InfinityModifier-n2.InfinityModifier == -2,
		Valid:            true,
	}
}

// Compares two numerics and returns:
//
//	-2 if any n1 or n2 is NaN
//	-1 if n1 <  n2
//	 0 if n1 == n2
//	+1 if n1 >  n2
//
// Panics if either n1 or n2 is not valid.
func Cmp(n1, n2 pgtype.Numeric) int {
	ensureValid(n1)
	ensureValid(n2)

	if n1.NaN || n2.NaN {
		return -2
	}

	if n1.InfinityModifier != pgtype.Finite || n2.InfinityModifier != pgtype.Finite {
		if n1.InfinityModifier < n2.InfinityModifier {
			return -1
		} else if n1.InfinityModifier > n2.InfinityModifier {
			return 1
		}
		return 0
	}

	if n1.Exp == n2.Exp {
		return n1.Int.Cmp(n2.Int)
	}

	rn1, rn2 := RescalePair(n1, n2)
	return rn1.Int.Cmp(rn2.Int)
}

func IsPositive(n pgtype.Numeric) bool {
	return Cmp(n, Zero) > 0
}

// Returns min(n1, n2).
//
// Panics if either n1 or n2 is not valid.
func Min(n1, n2 pgtype.Numeric) pgtype.Numeric {
	ensureValid(n1)
	ensureValid(n2)

	cmp := Cmp(n1, n2)
	if cmp == -2 {
		return NaN
	} else if cmp >= 0 {
		return n2
	}
	return n1
}

func ensureValid(n pgtype.Numeric) {
	if !n.Valid {
		log.Panicf("%+v is not valid", n)
	}
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
