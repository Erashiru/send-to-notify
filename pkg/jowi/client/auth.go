package client

import (
	"crypto/sha256"
	"encoding/hex"
)

func generateSig(apiKey, apiSecret string) string {
	hash := sha256.Sum256([]byte(apiKey + apiSecret))

	result := hex.EncodeToString(hash[:])
	return result[0:10] + result[len(result)-5:]
}
