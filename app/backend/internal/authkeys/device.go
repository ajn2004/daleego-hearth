package authkeys

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type DeviceAPIKey struct {
	Plaintext string
	Prefix    string
	Hash      string
}

func GenerateDeviceAPIKey() (DeviceAPIKey, error) {
	prefixBytes := make([]byte, 6)
	if _, err := rand.Read(prefixBytes); err != nil {
		return DeviceAPIKey{}, err
	}

	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return DeviceAPIKey{}, err
	}

	prefix := base64.RawURLEncoding.EncodeToString(prefixBytes)
	secret := base64.RawURLEncoding.EncodeToString(secretBytes)

	plaintext := fmt.Sprintf("hearth_dev_%s_%s", prefix, secret)

	sum := sha256.Sum256([]byte(plaintext))
	hash := hex.EncodeToString(sum[:])

	return DeviceAPIKey{
		Plaintext: plaintext,
		Prefix:    prefix,
		Hash:      hash,
	}, nil
}
