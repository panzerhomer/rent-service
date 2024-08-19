package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

func VerifyPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return err
}
