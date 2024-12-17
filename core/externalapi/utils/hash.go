package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/rs/zerolog/log"
)

func EncodeSecret(data string, appSecret string) (string, error) {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(appSecret))

	// Write Data to it
	_, err := h.Write([]byte(data))

	if err != nil {
		log.Trace().Err(err).Msg("Cant write data to HMAC")
		return "", err
	}

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	return sha, nil
}
