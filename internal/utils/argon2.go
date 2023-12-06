package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	format  = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	version = argon2.Version
	keyLen  = 32
	saltLen = 16
)

func Argon2IDHash(plain string) (string, error) {
	var time uint32 = 3
	var memory uint32 = 64 * 1024
	var threads uint8 = 4

	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(plain), salt, time, memory, threads, keyLen)

	return fmt.Sprintf(
			format,
			version,
			memory,
			time,
			threads,
			base64.RawStdEncoding.EncodeToString(salt),
			base64.RawStdEncoding.EncodeToString(hash),
		),
		nil
}

func Argon2IDVerify(plain, hash string) (bool, error) {
	var time uint32
	var memory uint32
	var threads uint8
	hashParts := strings.Split(hash, "$")

	_, err := fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])
	if err != nil {
		return false, err
	}

	hashToCompare := argon2.IDKey([]byte(plain), salt, time, memory, threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1, nil
}
