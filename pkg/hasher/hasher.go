package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashAPIKey(secret, apiKey string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(apiKey))
	return hex.EncodeToString(mac.Sum(nil))
}
