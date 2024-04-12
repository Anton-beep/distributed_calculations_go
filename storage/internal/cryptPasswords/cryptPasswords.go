package cryptPasswords

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func ComparePasswordWithHash(existing string, incoming string) error {
	return bcrypt.CompareHashAndPassword([]byte(existing), []byte(incoming))
}
