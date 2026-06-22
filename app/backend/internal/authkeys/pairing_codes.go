package authkeys

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func NormalizePairingCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func HashPairingCode(code string, secret string) string {
	normalized := NormalizePairingCode(code)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(normalized))

	return hex.EncodeToString(mac.Sum(nil))
}
