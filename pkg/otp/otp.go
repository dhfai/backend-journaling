package otp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func Generate() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func Hash(otp, pepper string) string {
	data := otp + pepper
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func Verify(otp, pepper, otpHash string) bool {
	computed := Hash(otp, pepper)
	return computed == otpHash
}
