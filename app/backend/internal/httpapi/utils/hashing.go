package httputil

import "golang.org/x/crypto/bcrypt"

func HashValue(valueToHash string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(valueToHash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyHashValue(value string, hashValue string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashValue), []byte(value))
	return err == nil
}
