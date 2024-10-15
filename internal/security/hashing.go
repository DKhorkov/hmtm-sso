package security

import "golang.org/x/crypto/bcrypt"

func Hash(value string, hashCost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), hashCost)
	return string(bytes), err
}

func ValidateHash(value, hashedValue string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	return err == nil
}
