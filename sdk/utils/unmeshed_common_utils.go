package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func CreateSecureHash(input string) (string, error) {
	hash := sha256.New()

	_, err := hash.Write([]byte(input))
	if err != nil {
		return "", fmt.Errorf("Error creating hash: %w", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
