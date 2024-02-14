package utils

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(src string) (pgtype.UUID, error) {
	var dst pgtype.UUID
	err := dst.Scan(src)
	return dst, err
}

// `EncodeUUID` converts a uuid byte array to UUID standard string form.
func EncodeUUID(src [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}
