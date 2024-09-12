package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(passwordByte []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
