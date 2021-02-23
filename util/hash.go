package util

import "golang.org/x/crypto/bcrypt"

func HashPin(rawString string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawString), 14)
	return string(hashed), err
}

func VerifyHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
