package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(src string) (pgtype.UUID, error) {
	var dst [16]byte

	switch len(src) {
	case 36:
		src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	case 32:
		// Dashes already stripped, assume valid
	default:
		// Assume invalid.
		return pgtype.UUID{}, fmt.Errorf("cannot parse UUID %v", src)
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return pgtype.UUID{}, err
	}

	copy(dst[:], buf)
	return pgtype.UUID{Bytes: dst, Valid: true}, err
}
